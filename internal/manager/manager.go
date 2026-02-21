package manager

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/mishamyrt/hapm/internal/git"
	"github.com/mishamyrt/hapm/internal/manifest"
	hapmpkg "github.com/mishamyrt/hapm/internal/package"
)

const maxApplyConcurrency = 20

type PackageManager struct {
	path     string
	lock     *Lockfile
	client   hapmpkg.GitClient
	registry hapmpkg.Registry
	packages map[string]hapmpkg.Package
	out      io.Writer
}

func New(path string, token string, out io.Writer) (*PackageManager, error) {
	if out == nil {
		out = os.Stdout
	}
	return NewWith(path, git.NewClient(token), hapmpkg.DefaultRegistry(), "_lock.json", out)
}

func NewWith(path string, client hapmpkg.GitClient, registry hapmpkg.Registry, lockfileName string, out io.Writer) (*PackageManager, error) {
	if lockfileName == "" {
		lockfileName = "_lock.json"
	}
	if out == nil {
		out = os.Stdout
	}
	manager := &PackageManager{
		path:     path,
		lock:     NewLockfile(filepath.Join(path, lockfileName)),
		client:   client,
		registry: registry,
		packages: map[string]hapmpkg.Package{},
		out:      out,
	}
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		if manager.lock.Exists() {
			if err := manager.bootFromLock(); err != nil {
				return nil, err
			}
		}
	} else {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return nil, err
		}
	}
	return manager, nil
}

func (m *PackageManager) SupportedTypes() []string {
	kinds := m.registry.SupportedKinds()
	sort.Strings(kinds)
	return kinds
}

func (m *PackageManager) GetVersions(location manifest.PackageLocation) ([]string, error) {
	return m.client.GetVersions(location.FullName)
}

func (m *PackageManager) bootFromLock() error {
	descriptions, err := m.lock.Load()
	if err != nil {
		return err
	}
	for _, description := range descriptions {
		constructor, ok := m.registry.Constructors[description.Kind]
		if !ok {
			return fmt.Errorf("unsupported package kind: %s", description.Kind)
		}
		pkg := constructor(description, m.path, m.client)
		m.packages[pkg.FullName()] = pkg
	}
	return nil
}

func (m *PackageManager) Diff(update []hapmpkg.PackageDescription, stableOnly bool) ([]PackageDiff, error) {
	updateFullNames := map[string]struct{}{}
	diffs := make([]PackageDiff, 0)

	for _, description := range update {
		current := description.Copy()
		if current.Version == "latest" {
			versions, err := m.client.GetVersions(current.FullName)
			if err != nil {
				return nil, err
			}
			current.Version = hapmpkg.FindLatestVersion(versions, stableOnly)
		}
		updateFullNames[current.FullName] = struct{}{}
		diff := PackageDiff{PackageDescription: current}
		if existing, ok := m.packages[current.FullName]; ok {
			if existing.Version() != current.Version {
				diff.CurrentVersion = existing.Version()
				diff.Operation = "switch"
			}
		} else {
			diff.Operation = "add"
		}
		if diff.Operation != "" {
			diffs = append(diffs, diff)
		}
	}

	for fullName, pkg := range m.packages {
		if _, ok := updateFullNames[fullName]; ok {
			continue
		}
		diff := PackageDiff{PackageDescription: pkg.Description(), Operation: "delete"}
		diffs = append(diffs, diff)
	}

	return diffs, nil
}

func (m *PackageManager) Apply(diffs []PackageDiff) error {
	type applyJob struct {
		index       int
		diff        PackageDiff
		pkg         hapmpkg.Package
		constructor hapmpkg.Constructor
	}
	type applyResult struct {
		index     int
		operation string
		fullName  string
		pkg       hapmpkg.Package
		err       error
	}

	jobs := make([]applyJob, 0, len(diffs))
	for i, diff := range diffs {
		switch diff.Operation {
		case "add":
			constructor, ok := m.registry.Constructors[diff.Kind]
			if !ok {
				return fmt.Errorf("unsupported package kind: %s", diff.Kind)
			}
			jobs = append(jobs, applyJob{index: i, diff: diff, constructor: constructor})
		case "delete":
			pkg, ok := m.packages[diff.FullName]
			if !ok {
				continue
			}
			jobs = append(jobs, applyJob{index: i, diff: diff, pkg: pkg})
		case "switch":
			pkg, ok := m.packages[diff.FullName]
			if !ok {
				return fmt.Errorf("package is not installed: %s", diff.FullName)
			}
			jobs = append(jobs, applyJob{index: i, diff: diff, pkg: pkg})
		default:
			return fmt.Errorf("unsupported operation: %s", diff.Operation)
		}
	}

	if len(jobs) == 0 {
		return m.lock.Dump(m.Descriptions())
	}

	workers := len(jobs)
	if workers > maxApplyConcurrency {
		workers = maxApplyConcurrency
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobCh := make(chan applyJob)
	resultCh := make(chan applyResult, len(jobs))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobCh:
				if !ok {
					return
				}

				result := applyResult{
					index:     job.index,
					operation: job.diff.Operation,
					fullName:  job.diff.FullName,
				}
				switch job.diff.Operation {
				case "add":
					pkg := job.constructor(job.diff.PackageDescription, m.path, m.client)
					if err := pkg.Setup(); err != nil {
						result.err = err
					} else {
						result.pkg = pkg
					}
				case "delete":
					result.err = job.pkg.Destroy()
				case "switch":
					result.err = job.pkg.Switch(job.diff.Version)
				}
				resultCh <- result
				if result.err != nil {
					cancel()
					return
				}
			}
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		defer close(jobCh)
		for _, job := range jobs {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case jobCh <- job:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make([]applyResult, 0, len(jobs))
	var firstErr error
	for result := range resultCh {
		results = append(results, result)
		if result.err != nil && firstErr == nil {
			firstErr = result.err
		}
	}
	if firstErr != nil {
		return firstErr
	}

	sort.Slice(results, func(i int, j int) bool {
		return results[i].index < results[j].index
	})
	for _, result := range results {
		switch result.operation {
		case "add":
			m.packages[result.fullName] = result.pkg
		case "delete":
			delete(m.packages, result.fullName)
		}
	}

	return m.lock.Dump(m.Descriptions())
}

func (m *PackageManager) Export(path string) error {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	kindsUsed := map[string]bool{}
	for _, pkg := range m.packages {
		if !kindsUsed[pkg.Kind()] {
			if hook, ok := m.registry.PreExport[pkg.Kind()]; ok {
				if err := hook(path); err != nil {
					return err
				}
			}
			kindsUsed[pkg.Kind()] = true
		}
		if err := pkg.Export(path); err != nil {
			return err
		}
	}
	kinds := make([]string, 0, len(kindsUsed))
	for kind := range kindsUsed {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	for _, kind := range kinds {
		if hook, ok := m.registry.PostExport[kind]; ok {
			if err := hook(path, m.out); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *PackageManager) Updates(stableOnly bool) ([]PackageDiff, error) {
	updates := make([]PackageDiff, 0)
	for _, pkg := range m.packages {
		latest, err := pkg.LatestVersion(stableOnly)
		if err != nil {
			return nil, err
		}
		latestVersion, err := hapmpkg.NewVersion(latest)
		if err != nil {
			continue
		}
		currentVersion, err := hapmpkg.NewVersion(pkg.Version())
		if err != nil {
			continue
		}
		if latestVersion.Compare(currentVersion) > 0 {
			updates = append(updates, PackageDiff{
				PackageDescription: hapmpkg.PackageDescription{FullName: pkg.FullName(), Kind: pkg.Kind(), Version: latest},
				CurrentVersion:     pkg.Version(),
				Operation:          "switch",
			})
		}
	}
	return updates, nil
}

func (m *PackageManager) Descriptions() []hapmpkg.PackageDescription {
	descriptions := make([]hapmpkg.PackageDescription, 0, len(m.packages))
	for _, pkg := range m.packages {
		descriptions = append(descriptions, pkg.Description())
	}
	return descriptions
}
