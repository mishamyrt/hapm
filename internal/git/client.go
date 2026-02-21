package git

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultAPIBaseURL = "https://api.github.com"
	defaultWebBaseURL = "https://github.com"
)

type Client struct {
	httpClient *http.Client
	apiBaseURL string
	webBaseURL string
	token      string
}

func NewClient(token string) *Client {
	apiBase := os.Getenv("HAPM_GITHUB_API_BASE_URL")
	if apiBase == "" {
		apiBase = defaultAPIBaseURL
	}
	webBase := os.Getenv("HAPM_GITHUB_WEB_BASE_URL")
	if webBase == "" {
		webBase = defaultWebBaseURL
	}
	return &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		apiBaseURL: strings.TrimRight(apiBase, "/"),
		webBaseURL: strings.TrimRight(webBase, "/"),
		token:      token,
	}
}

type tagItem struct {
	Name string `json:"name"`
}

type contentItem struct {
	Content string `json:"content"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type release struct {
	Assets []releaseAsset `json:"assets"`
}

func (c *Client) GetVersions(fullName string) ([]string, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/tags", c.apiBaseURL, fullName)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var tags []tagItem
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, err
	}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag.Name)
	}
	return result, nil
}

func (c *Client) GetTreeFile(fullName string, branch string, filePath string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/contents/%s?ref=%s", c.apiBaseURL, fullName, url.PathEscape(filePath), url.QueryEscape(branch))
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var content contentItem
	if err := json.Unmarshal(body, &content); err != nil {
		return nil, err
	}
	if content.Content == "" {
		return nil, fmt.Errorf("content is empty")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(content.Content, "\n", ""))
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (c *Client) GetReleaseFile(fullName string, branch string, filename string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/releases/tags/%s", c.apiBaseURL, fullName, url.QueryEscape(branch))
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var rel release
	if err := json.Unmarshal(body, &rel); err != nil {
		return nil, err
	}
	for _, asset := range rel.Assets {
		if asset.Name == filename {
			return c.get(asset.BrowserDownloadURL)
		}
	}
	return nil, fmt.Errorf("asset %s not found", filename)
}

func (c *Client) GetTarball(fullName string, branch string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s/tarball/%s", c.webBaseURL, fullName, url.PathEscape(branch))
	return c.get(endpoint)
}

func RepoURL(fullName string) string {
	return fmt.Sprintf("%s/%s", defaultWebBaseURL, fullName)
}

func (c *Client) get(endpoint string) ([]byte, error) {
	if parsed, err := url.Parse(endpoint); err == nil && parsed.Scheme == "file" {
		return os.ReadFile(parsed.Path)
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) == 0 {
			return nil, fmt.Errorf("http status: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("http status: %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return io.ReadAll(resp.Body)
}
