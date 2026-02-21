package manifest

import (
	"fmt"
	"os"
	"sort"

	"github.com/mishamyrt/hapm/internal/hapkg"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Path      string
	Values    []hapkg.PackageDescription
	HasLatest []string
}

func New(path string) *Manifest {
	return &Manifest{Path: path, Values: make([]hapkg.PackageDescription, 0), HasLatest: make([]string, 0)}
}

func (m *Manifest) Set(fullName string, version string, kind string) error {
	for i := range m.Values {
		if m.Values[i].FullName == fullName {
			m.Values[i].Version = version
			return nil
		}
	}
	if kind == "" {
		return fmt.Errorf("package type is not declared")
	}
	m.Values = append(m.Values, hapkg.PackageDescription{FullName: fullName, Version: version, Kind: kind})
	return nil
}

func (m *Manifest) Init(types []string) error {
	template := map[string][]string{}
	for _, packageType := range types {
		template[packageType] = []string{}
	}
	content, err := yaml.Marshal(template)
	if err != nil {
		return err
	}
	return os.WriteFile(m.Path, content, 0o644)
}

func (m *Manifest) Dump() error {
	content := map[string][]string{}
	for _, pkg := range m.Values {
		location := pkg.FullName + "@" + pkg.Version
		content[pkg.Kind] = append(content[pkg.Kind], location)
	}
	data, err := yaml.Marshal(content)
	if err != nil {
		return err
	}
	return os.WriteFile(m.Path, data, 0o644)
}

func (m *Manifest) Load() error {
	stream, err := os.ReadFile(m.Path)
	if err != nil {
		return err
	}
	if len(stream) == 0 {
		return fmt.Errorf("manifest is empty")
	}
	raw := map[string]any{}
	if err := yaml.Unmarshal(stream, &raw); err != nil {
		return err
	}
	if raw == nil {
		return fmt.Errorf("manifest is empty")
	}
	m.Values = m.Values[:0]
	m.HasLatest = m.HasLatest[:0]
	keys := make([]string, 0, len(raw))
	for key := range raw {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		values, err := ParseCategory(raw, key)
		if err != nil {
			return err
		}
		for _, value := range values {
			if value.Version == "latest" {
				m.HasLatest = append(m.HasLatest, value.FullName)
			}
		}
		m.Values = append(m.Values, values...)
	}
	return nil
}
