package hapm

import (
	"io"
	"os"

	"github.com/mishamyrt/hapm/internal/report"
)

const (
	storageDir   = ".hapm"
	manifestPath = "hapm.yaml"
	tokenVar     = "GITHUB_PAT"
)

// GlobalOptions describe global CLI flags.
type GlobalOptions struct {
	Manifest string
	Storage  string
	Dry      bool
}

// SyncOptions describe sync command options.
type SyncOptions struct {
	AllowUnstable bool
}

// InstallOptions describe install command options.
type InstallOptions struct {
	Entries       []string
	PackageType   string
	AllowUnstable bool
}

// UpdatesOptions describe updates command options.
type UpdatesOptions struct {
	AllowUnstable bool
}

// App coordinates command business logic.
type App struct {
	out      io.Writer
	errOut   io.Writer
	reporter report.Reporter
	globals  GlobalOptions
}

// DefaultGlobalOptions returns default values for global flags.
func DefaultGlobalOptions() GlobalOptions {
	return GlobalOptions{
		Manifest: manifestPath,
		Storage:  storageDir,
		Dry:      false,
	}
}

// New creates application instance.
func New(out io.Writer, errOut io.Writer) *App {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return &App{
		out:      out,
		errOut:   errOut,
		reporter: report.New(out),
		globals:  DefaultGlobalOptions(),
	}
}

// SetGlobals sets current global options.
func (a *App) SetGlobals(opts GlobalOptions) {
	if opts.Manifest == "" {
		opts.Manifest = manifestPath
	}
	if opts.Storage == "" {
		opts.Storage = storageDir
	}
	a.globals = opts
}

func (a *App) token() string {
	return os.Getenv(tokenVar)
}

func (a *App) warnNoToken() {
	if a.token() == "" {
		a.reporter.NoToken(tokenVar)
	}
}
