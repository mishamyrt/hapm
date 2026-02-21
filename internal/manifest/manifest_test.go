package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		in       string
		fullName string
		version  string
		ok       bool
	}{
		{"mishamyrt/myrt_desk_hass@v0.2.4", "mishamyrt/myrt_desk_hass", "v0.2.4", true},
		{"mishamyrt/myrt_desk_hass", "mishamyrt/myrt_desk_hass", "latest", true},
		{"https://github.com/mishamyrt/myrt_desk_hass", "mishamyrt/myrt_desk_hass", "latest", true},
		{"https://github.com/mishamyrt/myrt_desk_hass/releases/tag/v0.2.4", "mishamyrt/myrt_desk_hass", "v0.2.4", true},
		{"github.com/mishamyrt/myrt_desk_hass", "mishamyrt/myrt_desk_hass", "latest", true},
		{"hello", "", "", false},
	}
	for _, tc := range tests {
		location, ok := ParseLocation(tc.in)
		if ok != tc.ok {
			t.Fatalf("unexpected parse status for %q", tc.in)
		}
		if !tc.ok {
			continue
		}
		if location.FullName != tc.fullName || location.Version != tc.version {
			t.Fatalf("unexpected parse result for %q: %+v", tc.in, location)
		}
	}
}

func TestParseCategory(t *testing.T) {
	manifestContent := map[string]any{"integrations": []any{"foo/bar@v1.0.0"}}
	items, err := ParseCategory(manifestContent, "integrations")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one item")
	}
	if items[0].Kind != "integrations" {
		t.Fatalf("unexpected kind")
	}
}

func TestManifestInitLoadDumpSet(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "hapm.yaml")
	manifest := New(path)
	if err := manifest.Init([]string{"integrations", "plugins"}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "integrations") || !strings.Contains(string(raw), "plugins") {
		t.Fatalf("init template is wrong: %s", string(raw))
	}

	if err := manifest.Set("foo/bar", "v1.0.0", "integrations"); err != nil {
		t.Fatal(err)
	}
	if err := manifest.Set("foo/bar", "v1.0.1", "integrations"); err != nil {
		t.Fatal(err)
	}
	if err := manifest.Set("foo/plugin", "latest", "plugins"); err != nil {
		t.Fatal(err)
	}
	if err := manifest.Dump(); err != nil {
		t.Fatal(err)
	}

	loaded := New(path)
	if err := loaded.Load(); err != nil {
		t.Fatal(err)
	}
	if len(loaded.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(loaded.Values))
	}
	if len(loaded.HasLatest) != 1 || loaded.HasLatest[0] != "foo/plugin" {
		t.Fatalf("unexpected latest list: %+v", loaded.HasLatest)
	}
}

func TestManifestSetRequiresKind(t *testing.T) {
	manifest := New("unused")
	if err := manifest.Set("foo/bar", "v1.0.0", ""); err == nil {
		t.Fatalf("expected error")
	}
}

func TestManifestLoadErrors(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "hapm.yaml")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := New(path)
	if err := manifest.Load(); err == nil {
		t.Fatalf("expected empty manifest error")
	}
}
