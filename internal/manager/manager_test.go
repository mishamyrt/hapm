package manager

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
}

func (p *fakePackage) Description() hapmpkg.PackageDescription { return p.desc }
func (p *fakePackage) FullName() string                        { return p.desc.FullName }
func (p *fakePackage) Version() string                         { return p.desc.Version }
func (p *fakePackage) Kind() string                            { return p.desc.Kind }

func (p *fakePackage) Setup() error {
	return os.WriteFile(p.filePath(p.desc.Version), []byte(p.desc.Version), 0o644)
}

func (p *fakePackage) Switch(version string) error {
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
