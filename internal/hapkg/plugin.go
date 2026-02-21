package hapkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const PluginKind = "plugins"

const pluginFolderName = "www/custom_lovelace"

type PluginPackage struct {
	base BasePackage
}

func NewPluginPackage(description PackageDescription, rootPath string, client GitClient) Package {
	return &PluginPackage{base: newBasePackage(description, rootPath, "js", PluginKind, client)}
}

func (p *PluginPackage) Description() PackageDescription { return p.base.Description() }
func (p *PluginPackage) FullName() string                { return p.base.FullName() }
func (p *PluginPackage) Version() string                 { return p.base.Version() }
func (p *PluginPackage) Kind() string                    { return p.base.Kind() }
func (p *PluginPackage) Destroy() error                  { return p.base.Destroy() }
func (p *PluginPackage) LatestVersion(stableOnly bool) (string, error) {
	return p.base.LatestVersion(stableOnly)
}

func (p *PluginPackage) Setup() error {
	if p.base.version == "latest" {
		return fmt.Errorf("version is unknown")
	}
	return p.downloadScript(p.base.version)
}

func (p *PluginPackage) Switch(version string) error {
	if err := p.downloadScript(version); err != nil {
		return err
	}
	if err := os.Remove(p.base.Path("")); err != nil {
		return err
	}
	p.base.version = version
	return nil
}

func (p *PluginPackage) Export(path string) error {
	target := filepath.Join(path, pluginFolderName, p.base.name+".js")
	content, err := os.ReadFile(p.base.Path(""))
	if err != nil {
		return err
	}
	return os.WriteFile(target, content, 0o644)
}

func PluginPreExport(path string) error {
	return os.MkdirAll(filepath.Join(path, pluginFolderName), 0o755)
}

func PluginPostExport(path string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(path, pluginFolderName))
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}

func (p *PluginPackage) downloadScript(version string) error {
	content, err := p.getScript(version)
	if err != nil {
		return err
	}
	return os.WriteFile(p.base.Path(version), content, 0o644)
}

func (p *PluginPackage) getScript(version string) ([]byte, error) {
	pluginName := strings.TrimPrefix(p.base.name, "lovelace-")
	pluginFiles := []string{
		pluginName + ".js",
		pluginName + "-bundle.js",
	}
	for _, pluginFile := range pluginFiles {
		content, err := p.base.client.GetTreeFile(p.base.fullName, version, "dist/"+pluginFile)
		if err == nil && len(content) > 0 {
			return content, nil
		}
		content, err = p.base.client.GetTreeFile(p.base.fullName, version, pluginFile)
		if err == nil && len(content) > 0 {
			return content, nil
		}
		content, err = p.base.client.GetReleaseFile(p.base.fullName, version, pluginFile)
		if err == nil && len(content) > 0 {
			return content, nil
		}
	}
	return nil, fmt.Errorf("plugin script is not found: %s@%s", p.base.fullName, version)
}
