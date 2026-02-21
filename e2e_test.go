//go:build e2e
// +build e2e

package main_test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var (
	buildOnce sync.Once
	buildErr  error
	binPath   string
)

func TestScripts(t *testing.T) {
	hapmBin := buildBinary(t)
	params := testscript.Params{
		Dir: filepath.Join("testdata", "scripts"),
		Setup: func(env *testscript.Env) error {
			env.Setenv("PATH", filepath.Dir(hapmBin)+string(os.PathListSeparator)+os.Getenv("PATH"))
			env.Setenv("HAPM_DISABLE_PROGRESS", "1")
			if err := seedFixtures(env.WorkDir); err != nil {
				return err
			}
			return seedFakeAPI(env.WorkDir)
		},
	}
	testscript.Run(t, params)
}

func buildBinary(t *testing.T) string {
	t.Helper()
	buildOnce.Do(func() {
		tmp, err := os.MkdirTemp("", "hapm-go-bin")
		if err != nil {
			buildErr = err
			return
		}
		name := "hapm"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		binPath = filepath.Join(tmp, name)
		cmd := exec.Command("go", "build", "-o", binPath, "./hapm.go")
		cmd.Dir = ".."
		cmd.Env = append(os.Environ(),
			"GOSUMDB=off",
			"GOPROXY=off",
			"GOCACHE="+filepath.Join(tmp, ".gocache"),
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			buildErr = &buildFailure{err: err, output: string(output)}
		}
	})
	if buildErr != nil {
		t.Fatalf("build hapm binary: %v", buildErr)
	}
	return binPath
}

type buildFailure struct {
	err    error
	output string
}

func (e *buildFailure) Error() string {
	return e.err.Error() + ": " + e.output
}

func seedFixtures(workDir string) error {
	source := filepath.Clean(filepath.Join("..", "testdata"))
	target := filepath.Join(workDir, "fixtures")
	if err := copyDir(source, target); err != nil {
		return err
	}
	return nil
}

func seedFakeAPI(workDir string) error {
	apiRoot := filepath.Join(workDir, "fake", "api", "repos", "foo", "bar")
	if err := os.MkdirAll(apiRoot, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(apiRoot, "tags"), []byte(`[{"name":"v1.0.0"},{"name":"v1.1.0"}]`), 0o644); err != nil {
		return err
	}
	webRoot := filepath.Join(workDir, "fake", "web")
	return os.MkdirAll(webRoot, 0o755)
}

func copyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
