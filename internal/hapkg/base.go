package hapkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GitClient interface {
	GetVersions(fullName string) ([]string, error)
	GetTreeFile(fullName string, branch string, filePath string) ([]byte, error)
	GetReleaseFile(fullName string, branch string, filename string) ([]byte, error)
	GetTarball(fullName string, branch string) ([]byte, error)
}

type Package interface {
	Description() PackageDescription
	FullName() string
	Version() string
	Kind() string
	Setup() error
	Switch(version string) error
	Destroy() error
	Export(path string) error
	LatestVersion(stableOnly bool) (string, error)
}

type BasePackage struct {
	kind      string
	extension string
	fullName  string
	version   string
	basePath  string
	name      string
	client    GitClient
}

func newBasePackage(description PackageDescription, rootPath string, extension string, kind string, client GitClient) BasePackage {
	return BasePackage{
		kind:      kind,
		extension: extension,
		fullName:  description.FullName,
		version:   description.Version,
		basePath:  filepath.Join(rootPath, strings.ReplaceAll(description.FullName, "/", "-")),
		name:      description.ShortName(),
		client:    client,
	}
}

func (b *BasePackage) Description() PackageDescription {
	return PackageDescription{FullName: b.fullName, Kind: b.kind, Version: b.version}
}

func (b *BasePackage) Path(version string) string {
	if version == "" {
		version = b.version
	}
	return fmt.Sprintf("%s@%s.%s", b.basePath, version, b.extension)
}

func (b *BasePackage) Destroy() error {
	return os.Remove(b.Path(""))
}

func (b *BasePackage) LatestVersion(stableOnly bool) (string, error) {
	versions, err := b.client.GetVersions(b.fullName)
	if err != nil {
		return "", err
	}
	return FindLatestVersion(versions, stableOnly), nil
}

func (b *BasePackage) FullName() string {
	return b.fullName
}

func (b *BasePackage) Version() string {
	return b.version
}

func (b *BasePackage) Kind() string {
	return b.kind
}
