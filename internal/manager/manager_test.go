package manager

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	hapmpkg "github.com/mishamyrt/hapm/internal/package"
)

type fakeClient struct {
	versions map[string][]string
}

func (f fakeClient) GetVersions(fullName string) ([]string, error) {
	if versions, ok := f.versions[fullName]; ok {
		return versions, nil
	}
	return nil, errors.New("versions not found")
}

func (f fakeClient) GetTreeFile(string, string, string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (f fakeClient) GetReleaseFile(string, string, string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (f fakeClient) GetTarball(string, string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

type fakePackage struct {
	desc   hapmpkg.PackageDescription
	root   string
	latest string

	setupFn   func(*fakePackage) error
	switchFn  func(*fakePackage, string) error
	destroyFn func(*fakePackage) error
}

func (p *fakePackage) Description() hapmpkg.PackageDescription { return p.desc }
func (p *fakePackage) FullName() string                        { return p.desc.FullName }
func (p *fakePackage) Version() string                         { return p.desc.Version }
func (p *fakePackage) Kind() string                            { return p.desc.Kind }

func (p *fakePackage) Setup() error {
	if p.setupFn != nil {
		return p.setupFn(p)
	}
	return os.WriteFile(p.filePath(p.desc.Version), []byte(p.desc.Version), 0o644)
}

func (p *fakePackage) Switch(version string) error {
	if p.switchFn != nil {
		return p.switchFn(p, version)
	}
	if err := os.WriteFile(p.filePath(version), []byte(version), 0o644); err != nil {
		return err
	}
	if err := os.Remove(p.filePath(p.desc.Version)); err != nil {
		return err
	}
	p.desc.Version = version
	return nil
}

func (p *fakePackage) Destroy() error {
	if p.destroyFn != nil {
		return p.destroyFn(p)
	}
	return os.Remove(p.filePath(p.desc.Version))
}

func (p *fakePackage) Export(path string) error {
	if err := os.MkdirAll(filepath.Join(path, p.desc.Kind), 0o755); err != nil {
		return err
	}
	name := strings.ReplaceAll(p.desc.FullName, "/", "-") + ".txt"
	return os.WriteFile(filepath.Join(path, p.desc.Kind, name), []byte(p.desc.Version), 0o644)
}

func (p *fakePackage) LatestVersion(bool) (string, error) {
	return p.latest, nil
}

func (p *fakePackage) filePath(version string) string {
	name := strings.ReplaceAll(p.desc.FullName, "/", "-")
	return filepath.Join(p.root, name+"@"+version+".pkg")
}

func TestLockfileRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	lock := NewLockfile(filepath.Join(tmp, "_lock.json"))
	descriptions := []hapmpkg.PackageDescription{{FullName: "foo/bar", Version: "v1.0.0", Kind: "integrations"}}
	if err := lock.Dump(descriptions); err != nil {
		t.Fatal(err)
	}
	items, err := lock.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].FullName != "foo/bar" {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestManagerDiffApplyUpdatesAndExport(t *testing.T) {
	tmp := t.TempDir()
	out := &bytes.Buffer{}
	latestByName := map[string]string{"foo/bar": "v2.0.0"}
	registry := hapmpkg.Registry{
		Constructors: map[string]hapmpkg.Constructor{
			"integrations": func(description hapmpkg.PackageDescription, rootPath string, _ hapmpkg.GitClient) hapmpkg.Package {
				return &fakePackage{desc: description, root: rootPath, latest: latestByName[description.FullName]}
			},
		},
		PreExport: map[string]func(path string) error{
			"integrations": func(path string) error {
				return os.MkdirAll(filepath.Join(path, "integrations"), 0o755)
			},
		},
		PostExport: map[string]func(path string, out io.Writer) error{
			"integrations": func(_ string, out io.Writer) error {
				_, _ = out.Write([]byte("post-export\n"))
				return nil
			},
		},
	}

	manager, err := NewWith(tmp, fakeClient{versions: map[string][]string{"foo/bar": {"v1.0.0", "v1.2.0"}}}, registry, "_lock.json", out)
	if err != nil {
		t.Fatal(err)
	}

	update := []hapmpkg.PackageDescription{{FullName: "foo/bar", Version: "latest", Kind: "integrations"}}
	diff, err := manager.Diff(update, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(diff) != 1 || diff[0].Operation != "add" || diff[0].Version != "v1.2.0" {
		t.Fatalf("unexpected diff: %+v", diff)
	}

	if err := manager.Apply(diff); err != nil {
		t.Fatal(err)
	}
	descriptions := manager.Descriptions()
	if len(descriptions) != 1 || descriptions[0].Version != "v1.2.0" {
		t.Fatalf("unexpected descriptions: %+v", descriptions)
	}

	updates, err := manager.Updates(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(updates) != 1 || updates[0].Version != "v2.0.0" {
		t.Fatalf("unexpected updates: %+v", updates)
	}

	exportPath := filepath.Join(tmp, "export")
	if err := manager.Export(exportPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(exportPath, "integrations", "foo-bar.txt")); err != nil {
		t.Fatalf("missing export file: %v", err)
	}
	if !strings.Contains(out.String(), "post-export") {
		t.Fatalf("missing post export output: %s", out.String())
	}

	deleteDiff, err := manager.Diff([]hapmpkg.PackageDescription{}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(deleteDiff) != 1 || deleteDiff[0].Operation != "delete" {
		t.Fatalf("unexpected delete diff: %+v", deleteDiff)
	}
}

func TestManagerApplyLimitsConcurrencyTo20(t *testing.T) {
	tmp := t.TempDir()
	started := make(chan struct{}, 64)
	release := make(chan struct{})
	var mu sync.Mutex
	active := 0
	maxActive := 0

	registry := hapmpkg.Registry{
		Constructors: map[string]hapmpkg.Constructor{
			"integrations": func(description hapmpkg.PackageDescription, rootPath string, _ hapmpkg.GitClient) hapmpkg.Package {
				return &fakePackage{
					desc: description,
					root: rootPath,
					setupFn: func(_ *fakePackage) error {
						mu.Lock()
						active++
						if active > maxActive {
							maxActive = active
						}
						mu.Unlock()

						started <- struct{}{}
						<-release

						mu.Lock()
						active--
						mu.Unlock()
						return nil
					},
				}
			},
		},
	}

	manager, err := NewWith(tmp, fakeClient{}, registry, "_lock.json", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	diffs := make([]PackageDiff, 0, 30)
	for i := 0; i < 30; i++ {
		diffs = append(diffs, PackageDiff{
			PackageDescription: hapmpkg.PackageDescription{
				FullName: fmt.Sprintf("foo/pkg-%02d", i),
				Kind:     "integrations",
				Version:  "v1.0.0",
			},
			Operation: "add",
		})
	}

	applyErr := make(chan error, 1)
	go func() {
		applyErr <- manager.Apply(diffs)
	}()

	for i := 0; i < maxApplyConcurrency; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for started task #%d", i+1)
		}
	}

	select {
	case <-started:
		t.Fatal("expected at most 20 concurrent tasks before release")
	case <-time.After(150 * time.Millisecond):
	}

	close(release)

	select {
	case err := <-applyErr:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for apply completion")
	}

	mu.Lock()
	gotMax := maxActive
	gotActive := active
	mu.Unlock()
	if gotMax != maxApplyConcurrency {
		t.Fatalf("unexpected max concurrency: got %d, want %d", gotMax, maxApplyConcurrency)
	}
	if gotActive != 0 {
		t.Fatalf("unexpected active tasks after apply: %d", gotActive)
	}
	if len(manager.Descriptions()) != 30 {
		t.Fatalf("unexpected package count: %d", len(manager.Descriptions()))
	}
}

func TestManagerApplyCancelOnFirstError(t *testing.T) {
	tmp := t.TempDir()
	sentinelErr := errors.New("boom")
	release := make(chan struct{})
	var releaseOnce sync.Once
	var started atomic.Int32

	registry := hapmpkg.Registry{
		Constructors: map[string]hapmpkg.Constructor{
			"integrations": func(description hapmpkg.PackageDescription, rootPath string, _ hapmpkg.GitClient) hapmpkg.Package {
				return &fakePackage{
					desc: description,
					root: rootPath,
					setupFn: func(p *fakePackage) error {
						started.Add(1)
						if p.desc.FullName == "foo/pkg-00" {
							releaseOnce.Do(func() {
								close(release)
							})
							return sentinelErr
						}
						<-release
						return nil
					},
				}
			},
		},
	}

	manager, err := NewWith(tmp, fakeClient{}, registry, "_lock.json", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	lockPath := filepath.Join(tmp, "_lock.json")
	originalLock := []byte(`[{"full_name":"seed/pkg","kind":"integrations","version":"v1.0.0"}]`)
	if err := os.WriteFile(lockPath, originalLock, 0o644); err != nil {
		t.Fatal(err)
	}

	diffs := make([]PackageDiff, 0, 30)
	for i := 0; i < 30; i++ {
		diffs = append(diffs, PackageDiff{
			PackageDescription: hapmpkg.PackageDescription{
				FullName: fmt.Sprintf("foo/pkg-%02d", i),
				Kind:     "integrations",
				Version:  "v1.0.0",
			},
			Operation: "add",
		})
	}

	err = manager.Apply(diffs)
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := int(started.Load()); got >= len(diffs) {
		t.Fatalf("expected cancellation before all jobs started, got %d started of %d", got, len(diffs))
	}

	after, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, originalLock) {
		t.Fatalf("lockfile was unexpectedly changed: %s", string(after))
	}
}
