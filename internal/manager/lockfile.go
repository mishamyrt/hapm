package manager

import (
	"encoding/json"
	"os"

	hapmpkg "github.com/mishamyrt/hapm/internal/package"
)

type Lockfile struct {
	path string
}

func NewLockfile(path string) *Lockfile {
	return &Lockfile{path: path}
}

func (l *Lockfile) Exists() bool {
	_, err := os.Stat(l.path)
	return err == nil
}

func (l *Lockfile) Dump(descriptions []hapmpkg.PackageDescription) error {
	content, err := json.Marshal(descriptions)
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, content, 0o644)
}

func (l *Lockfile) Load() ([]hapmpkg.PackageDescription, error) {
	content, err := os.ReadFile(l.path)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return []hapmpkg.PackageDescription{}, nil
	}
	items := []hapmpkg.PackageDescription{}
	if err := json.Unmarshal(content, &items); err != nil {
		return nil, err
	}
	return items, nil
}
