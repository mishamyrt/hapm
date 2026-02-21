package hapkg

import "strings"

type PackageDescription struct {
	FullName string `json:"full_name" yaml:"full_name"`
	Version  string `json:"version" yaml:"version"`
	Kind     string `json:"kind" yaml:"kind"`
}

func (d PackageDescription) Copy() PackageDescription {
	return PackageDescription{
		FullName: d.FullName,
		Version:  d.Version,
		Kind:     d.Kind,
	}
}

func (d PackageDescription) ShortName() string {
	i := strings.LastIndex(d.FullName, "/")
	if i < 0 {
		return d.FullName
	}
	return d.FullName[i+1:]
}
