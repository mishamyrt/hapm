package manager

import hapmpkg "github.com/mishamyrt/hapm/internal/package"

type PackageDiff struct {
	hapmpkg.PackageDescription
	Operation      string `json:"operation"`
	CurrentVersion string `json:"current_version,omitempty"`
}
