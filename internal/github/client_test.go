package github

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestClientGetVersionsAndTarball(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.URL.String() {
			case "https://api.local/repos/foo/bar/tags":
				return newResponse(200, `[{"name":"v1.0.0"},{"name":"v1.0.1"}]`), nil
			case "https://web.local/foo/bar/tarball/v1.0.1":
				return newResponse(200, "tarball"), nil
			default:
				return newResponse(404, "not found"), nil
			}
		})},
		apiBaseURL: "https://api.local",
		webBaseURL: "https://web.local",
		token:      "token",
	}

	versions, err := client.GetVersions("foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if len(versions) != 2 || versions[1] != "v1.0.1" {
		t.Fatalf("unexpected versions: %+v", versions)
	}

	content, err := client.GetTarball("foo/bar", "v1.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "tarball" {
		t.Fatalf("unexpected tarball payload: %q", string(content))
	}
}

func TestClientTreeAndReleaseFallback(t *testing.T) {
	script := []byte("console.log('ok')")
	encoded := base64.StdEncoding.EncodeToString(script)
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			url := req.URL.String()
			switch {
			case strings.HasPrefix(url, "https://api.local/repos/foo/bar/contents/dist%2Fplugin.js"):
				return newResponse(200, fmt.Sprintf(`{"content":%q}`, encoded)), nil
			case url == "https://api.local/repos/foo/bar/releases/tags/v1.0.0":
				return newResponse(200, `{"assets":[{"name":"plugin.js","browser_download_url":"https://api.local/asset/plugin.js"}]}`), nil
			case url == "https://api.local/asset/plugin.js":
				return newResponse(200, string(script)), nil
			default:
				return newResponse(404, "not found"), nil
			}
		})},
		apiBaseURL: "https://api.local",
		webBaseURL: "https://web.local",
	}

	content, err := client.GetTreeFile("foo/bar", "v1.0.0", "dist/plugin.js")
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(script) {
		t.Fatalf("unexpected tree content: %q", string(content))
	}

	releaseContent, err := client.GetReleaseFile("foo/bar", "v1.0.0", "plugin.js")
	if err != nil {
		t.Fatal(err)
	}
	if string(releaseContent) != string(script) {
		t.Fatalf("unexpected release content: %q", string(releaseContent))
	}
}

func TestRepoURL(t *testing.T) {
	if got := RepoURL("foo/bar"); got != "https://github.com/foo/bar" {
		t.Fatalf("unexpected repo url: %s", got)
	}
}

func TestMain(m *testing.M) {
	_ = os.Unsetenv("HAPM_GITHUB_API_BASE_URL")
	_ = os.Unsetenv("HAPM_GITHUB_WEB_BASE_URL")
	os.Exit(m.Run())
}
