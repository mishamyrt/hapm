package manager

import "github.com/mishamyrt/hapm/internal/hapkg"

type PackageDiff struct {
	hapkg.PackageDescription
	Operation      string `json:"operation"`
	CurrentVersion string `json:"current_version,omitempty"`
}
