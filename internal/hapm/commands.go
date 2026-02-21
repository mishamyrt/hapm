package hapm

import (
	"fmt"

	"github.com/mishamyrt/hapm/internal/manager"
	"github.com/mishamyrt/hapm/internal/manifest"
	"github.com/mishamyrt/hapm/internal/report"
)

func (a *App) newManager() (*manager.PackageManager, error) {
	store, err := manager.New(a.globals.Storage, a.token())
	if err != nil {
		return nil, a.handledError("creating package manager", err)
	}
	return store, nil
}

// Init creates empty manifest from supported package kinds.
func (a *App) Init() error {
	store, err := a.newManager()
	if err != nil {
		return err
	}
	manifestFile := manifest.New(a.globals.Manifest)
	if err := manifestFile.Init(store.SupportedTypes()); err != nil {
		return a.handledError("initializing manifest", err)
	}
	return nil
}

// Sync synchronizes storage with current manifest.
func (a *App) Sync(opts SyncOptions) error {
	store, err := a.newManager()
	if err != nil {
		return err
	}
	return a.synchronize(store, !opts.AllowUnstable, nil)
}

// Install adds package locations to manifest and synchronizes.
func (a *App) Install(opts InstallOptions) error {
	store, err := a.newManager()
	if err != nil {
		return err
	}
	if len(opts.Entries) == 0 {
		return a.handledMessage("install requires at least one package location")
	}

	manifestFile := manifest.New(a.globals.Manifest)
	if err := manifestFile.Load(); err != nil {
		return a.handledError("parsing manifest", err)
	}

	for _, entry := range opts.Entries {
		location, ok := manifest.ParseLocation(entry)
		if !ok {
			a.reporter.Exception("installing package", fmt.Errorf("wrong location: %s", entry))
			if opts.PackageType == "" {
				a.reporter.Warning("--type parameter is not specified.\nThis option is required when installing new packages")
			}
			continue
		}
		if err := manifestFile.Set(location.FullName, location.Version, opts.PackageType); err != nil {
			a.reporter.Exception("installing package", err)
			if opts.PackageType == "" {
				a.reporter.Warning("--type parameter is not specified.\nThis option is required when installing new packages")
			}
		}
	}

	return a.synchronize(store, !opts.AllowUnstable, manifestFile)
}

// PrintUpdates prints list of packages that can be upgraded.
func (a *App) PrintUpdates(opts UpdatesOptions) error {
	store, err := a.newManager()
	if err != nil {
		return err
	}

	stableOnly := !opts.AllowUnstable
	if !stableOnly {
		a.reporter.Warning("Search includes unstable versions")
	}

	progress := report.NewProgress(a.reporter.Out())
	progress.Start("Looking for package updates")
	diff, err := store.Updates(stableOnly)
	progress.Stop()
	if err != nil {
		return a.handledError("looking for package updates", err)
	}
	if len(diff) == 0 {
		a.reporter.UpToDate()
		return nil
	}
	a.reporter.Diff(diff, true, true)
	return nil
}

// PrintVersions prints all available versions for a package location.
func (a *App) PrintVersions(entries []string) error {
	store, err := a.newManager()
	if err != nil {
		return err
	}
	if len(entries) != 1 {
		return a.handledMessage("versions requires exactly one package location")
	}
	location, ok := manifest.ParseLocation(entries[0])
	if !ok {
		a.reporter.WrongFormat(entries[0])
		return HandledError(fmt.Errorf("wrong location format: %s", entries[0]))
	}

	progress := report.NewProgress(a.reporter.Out())
	progress.Start("Looking for package versions")
	tags, err := store.GetVersions(*location)
	progress.Stop()
	if err != nil {
		return a.handledError("looking for package versions", err)
	}
	a.reporter.Versions(location.FullName, tags)
	return nil
}

// List prints installed packages grouped by kind.
func (a *App) List() error {
	store, err := a.newManager()
	if err != nil {
		return err
	}
	a.reporter.Packages(store.Descriptions())
	return nil
}

// Export copies installed packages to target directory.
func (a *App) Export(entries []string) error {
	store, err := a.newManager()
	if err != nil {
		return err
	}
	if len(entries) != 1 {
		return a.handledMessage("export requires output path")
	}
	result, err := store.Export(entries[0])
	if err != nil {
		return a.handledError("exporting packages", err)
	}
	if files, ok := result.PostExportFiles["plugins"]; ok {
		a.reporter.PluginExportHint(files)
	}
	return nil
}

func (a *App) synchronize(
	store *manager.PackageManager,
	stableOnly bool,
	loadedManifest *manifest.Manifest,
) error {
	if loadedManifest == nil {
		loadedManifest = manifest.New(a.globals.Manifest)
		if err := loadedManifest.Load(); err != nil {
			return a.handledError("parsing manifest", err)
		}
	}

	progress := report.NewProgress(a.reporter.Out())
	if len(loadedManifest.HasLatest) > 0 {
		a.reporter.Latest(loadedManifest.HasLatest)
		progress.Start("Search for the latest versions")
	}
	diff, err := store.Diff(loadedManifest.Values, stableOnly)
	if len(loadedManifest.HasLatest) > 0 {
		progress.Stop()
	}
	if err != nil {
		return a.handledError("calculating changes", err)
	}
	a.reporter.Diff(diff, false, false)
	if a.globals.Dry {
		return nil
	}
	if len(diff) > 0 {
		a.warnNoToken()
		progress := report.NewProgress(a.reporter.Out())
		progress.Start("Synchronizing the changes")
		if err := store.Apply(diff); err != nil {
			progress.Stop()
			return a.handledError("synchronizing the changes", err)
		}
		progress.Stop()
	}
	a.reporter.Summary(diff)
	return nil
}
