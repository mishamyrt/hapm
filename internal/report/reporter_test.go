package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mishamyrt/hapm/internal/manager"
	hapmpkg "github.com/mishamyrt/hapm/internal/package"
)

func TestReporterDiffAndSummary(t *testing.T) {
	out := &bytes.Buffer{}
	r := New(out)
	diffs := []manager.PackageDiff{
		{PackageDescription: hapmpkg.PackageDescription{FullName: "foo/new", Kind: "integrations", Version: "v1.0.0"}, Operation: "add"},
		{PackageDescription: hapmpkg.PackageDescription{FullName: "foo/old", Kind: "integrations", Version: "v2.0.0"}, Operation: "switch", CurrentVersion: "v1.0.0"},
		{PackageDescription: hapmpkg.PackageDescription{FullName: "foo/drop", Kind: "integrations", Version: "v1.0.0"}, Operation: "delete"},
	}
	r.Diff(diffs, false, false)
	r.Summary(diffs)

	text := out.String()
	for _, needle := range []string{"Integrations:", "+ new", "* old", "- drop", "Done:"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("missing %q in output: %s", needle, text)
		}
	}
}

func TestReporterWarnings(t *testing.T) {
	out := &bytes.Buffer{}
	r := New(out)
	r.NoToken("GITHUB_PAT")
	r.WrongFormat("bad-format")
	text := out.String()
	if !strings.Contains(text, "$GITHUB_PAT is not defined") {
		t.Fatalf("missing no-token warning")
	}
	if !strings.Contains(text, "Wrong location format") {
		t.Fatalf("missing wrong format warning")
	}
}
