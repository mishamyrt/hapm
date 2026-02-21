package hapkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type fakeGitClient struct {
	versions map[string][]string
	tarballs map[string][]byte
	tree     map[string][]byte
	release  map[string][]byte
}

func (f fakeGitClient) GetVersions(fullName string) ([]string, error) {
	if versions, ok := f.versions[fullName]; ok {
		return versions, nil
	}
	return nil, errors.New("versions not found")
}

func (f fakeGitClient) GetTreeFile(fullName string, branch string, filePath string) ([]byte, error) {
	key := fullName + "@" + branch + ":" + filePath
	if content, ok := f.tree[key]; ok {
		return content, nil
	}
	return nil, errors.New("tree file not found")
}

func (f fakeGitClient) GetReleaseFile(fullName string, branch string, filename string) ([]byte, error) {
	key := fullName + "@" + branch + ":" + filename
	if content, ok := f.release[key]; ok {
		return content, nil
	}
	return nil, errors.New("release file not found")
}

func (f fakeGitClient) GetTarball(fullName string, branch string) ([]byte, error) {
	key := fullName + "@" + branch
	if content, ok := f.tarballs[key]; ok {
		return content, nil
	}
	return nil, errors.New("tarball not found")
}

func TestIntegrationPackageExport(t *testing.T) {
	tmp := t.TempDir()
	tarball := makeTarball(t, map[string]string{
		"repo-abc/custom_components/demo/manifest.json": `{}`,
		"repo-abc/custom_components/demo/__init__.py":   "",
		"repo-abc/README.md":                            "hello",
	})
	client := fakeGitClient{tarballs: map[string][]byte{"foo/demo@v1.0.0": tarball}}
	desc := PackageDescription{FullName: "foo/demo", Version: "v1.0.0", Kind: IntegrationKind}

	pkg := NewIntegrationPackage(desc, tmp, client)
	if err := pkg.Setup(); err != nil {
		t.Fatal(err)
	}

	exportDir := filepath.Join(tmp, "export")
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := IntegrationPreExport(exportDir); err != nil {
		t.Fatal(err)
	}
	if err := pkg.Export(exportDir); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(exportDir, "custom_components", "demo", "manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("expected exported component: %v", err)
	}
}

func TestPluginPackageFallbackAndExport(t *testing.T) {
	tmp := t.TempDir()
	script := []byte("console.log('ok')")
	client := fakeGitClient{
		tree: map[string][]byte{},
		release: map[string][]byte{
			"foo/lovelace-demo@v1.0.0:demo.js": script,
		},
	}
	_, content := setupAndExportPlugin(t, tmp, client)
	if string(content) != string(script) {
		t.Fatalf("unexpected script content: %s", string(content))
	}
}

func TestPluginPackageBundleReleaseFallbackAndExport(t *testing.T) {
	tmp := t.TempDir()
	script := []byte("console.log('release-bundle')")
	client := fakeGitClient{
		tree: map[string][]byte{},
		release: map[string][]byte{
			"foo/lovelace-demo@v1.0.0:demo-bundle.js": script,
		},
	}

	outputPath, content := setupAndExportPlugin(t, tmp, client)
	if string(content) != string(script) {
		t.Fatalf("unexpected script content: %s", string(content))
	}
	bundleOutputPath := filepath.Join(filepath.Dir(outputPath), "lovelace-demo-bundle.js")
	if _, err := os.Stat(bundleOutputPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected bundle output file: %v", err)
	}
}

func TestPluginPackageBundleDistFallback(t *testing.T) {
	tmp := t.TempDir()
	script := []byte("console.log('dist-bundle')")
	client := fakeGitClient{
		tree: map[string][]byte{
			"foo/lovelace-demo@v1.0.0:dist/demo-bundle.js": script,
		},
		release: map[string][]byte{},
	}

	_, content := setupAndExportPlugin(t, tmp, client)
	if string(content) != string(script) {
		t.Fatalf("unexpected script content: %s", string(content))
	}
}

func TestPluginPackageBundleRootFallback(t *testing.T) {
	tmp := t.TempDir()
	script := []byte("console.log('root-bundle')")
	client := fakeGitClient{
		tree: map[string][]byte{
			"foo/lovelace-demo@v1.0.0:demo-bundle.js": script,
		},
		release: map[string][]byte{},
	}

	_, content := setupAndExportPlugin(t, tmp, client)
	if string(content) != string(script) {
		t.Fatalf("unexpected script content: %s", string(content))
	}
}

func TestPluginPackagePrefersRegularScriptOverBundle(t *testing.T) {
	tmp := t.TempDir()
	regularScript := []byte("console.log('regular')")
	bundleScript := []byte("console.log('bundle')")
	client := fakeGitClient{
		tree: map[string][]byte{
			"foo/lovelace-demo@v1.0.0:dist/demo.js":        regularScript,
			"foo/lovelace-demo@v1.0.0:dist/demo-bundle.js": bundleScript,
		},
		release: map[string][]byte{},
	}

	_, content := setupAndExportPlugin(t, tmp, client)
	if string(content) != string(regularScript) {
		t.Fatalf("unexpected script content: %s", string(content))
	}
}

func TestPluginPackageNotFound(t *testing.T) {
	tmp := t.TempDir()
	client := fakeGitClient{tree: map[string][]byte{}, release: map[string][]byte{}}
	desc := PackageDescription{FullName: "foo/bar", Version: "v1.0.0", Kind: PluginKind}
	pkg := NewPluginPackage(desc, tmp, client)
	if err := pkg.Setup(); err == nil {
		t.Fatalf("expected setup error")
	}
}

func setupAndExportPlugin(t *testing.T, tmp string, client fakeGitClient) (string, []byte) {
	t.Helper()

	desc := PackageDescription{FullName: "foo/lovelace-demo", Version: "v1.0.0", Kind: PluginKind}
	pkg := NewPluginPackage(desc, tmp, client)
	if err := pkg.Setup(); err != nil {
		t.Fatal(err)
	}
	exportDir := filepath.Join(tmp, "export")
	if err := PluginPreExport(exportDir); err != nil {
		t.Fatal(err)
	}
	if err := pkg.Export(exportDir); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(exportDir, "www", "custom_lovelace", "lovelace-demo.js")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	return outputPath, content
}

func makeTarball(t *testing.T, files map[string]string) []byte {
	t.Helper()
	buffer := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buffer)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		header := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(content))}
		if filepath.Ext(name) == "" {
			header.Typeflag = tar.TypeDir
			header.Mode = 0o755
			header.Size = 0
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("write header %s: %v", name, err)
		}
		if header.Typeflag != tar.TypeDir {
			if _, err := tw.Write([]byte(content)); err != nil {
				t.Fatalf("write content %s: %v", name, err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
