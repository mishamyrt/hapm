package hapmpkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const IntegrationKind = "integrations"
const integrationFolderName = "custom_components"

type IntegrationPackage struct {
	base BasePackage
}

func NewIntegrationPackage(description PackageDescription, rootPath string, client GitClient) Package {
	return &IntegrationPackage{base: newBasePackage(description, rootPath, "tar.gz", IntegrationKind, client)}
}

func (p *IntegrationPackage) Description() PackageDescription { return p.base.Description() }
func (p *IntegrationPackage) FullName() string                { return p.base.FullName() }
func (p *IntegrationPackage) Version() string                 { return p.base.Version() }
func (p *IntegrationPackage) Kind() string                    { return p.base.Kind() }
func (p *IntegrationPackage) Destroy() error                  { return p.base.Destroy() }
func (p *IntegrationPackage) LatestVersion(stableOnly bool) (string, error) {
	return p.base.LatestVersion(stableOnly)
}

func (p *IntegrationPackage) Setup() error {
	if p.base.version == "latest" {
		return fmt.Errorf("version is unknown")
	}
	return p.downloadTarball(p.base.version)
}

func (p *IntegrationPackage) Switch(version string) error {
	if err := p.downloadTarball(version); err != nil {
		return err
	}
	if err := os.Remove(p.base.Path("")); err != nil {
		return err
	}
	p.base.version = version
	return nil
}

func (p *IntegrationPackage) Export(dest string) error {
	file, err := os.Open(p.base.Path(""))
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := gz.Close()
		if err == nil {
			err = closeErr
		}
	}()

	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := filepath.ToSlash(header.Name)
		idx := strings.Index(name, "/"+integrationFolderName+"/")
		if idx < 0 {
			continue
		}
		rel := name[idx+1:]
		target := filepath.Join(dest, rel)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, reader); err != nil {
				_ = out.Close()
				return err
			}
			if err := out.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *IntegrationPackage) downloadTarball(version string) error {
	content, err := p.base.client.GetTarball(p.base.fullName, version)
	if err != nil {
		return err
	}
	return os.WriteFile(p.base.Path(version), content, 0o644)
}

func IntegrationPreExport(path string) error {
	return os.Mkdir(filepath.Join(path, integrationFolderName), 0o755)
}

func IntegrationPostExport(_ string, _ io.Writer) error {
	return nil
}
