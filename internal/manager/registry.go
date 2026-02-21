package manager

import "github.com/mishamyrt/hapm/internal/hapkg"

type Constructor func(description hapkg.PackageDescription, rootPath string, client hapkg.GitClient) hapkg.Package

type Registry struct {
	Constructors map[string]Constructor
	PreExport    map[string]func(path string) error
	PostExport   map[string]func(path string) ([]string, error)
}

func (r Registry) SupportedKinds() []string {
	kinds := make([]string, 0, len(r.Constructors))
	for kind := range r.Constructors {
		kinds = append(kinds, kind)
	}
	return kinds
}

func DefaultRegistry() Registry {
	return Registry{
		Constructors: map[string]Constructor{
			hapkg.IntegrationKind: func(d hapkg.PackageDescription, r string, c hapkg.GitClient) hapkg.Package {
				return hapkg.NewIntegrationPackage(d, r, c)
			},
			hapkg.PluginKind: func(d hapkg.PackageDescription, r string, c hapkg.GitClient) hapkg.Package {
				return hapkg.NewPluginPackage(d, r, c)
			},
		},
		PreExport: map[string]func(path string) error{
			hapkg.IntegrationKind: hapkg.IntegrationPreExport,
			hapkg.PluginKind:      hapkg.PluginPreExport,
		},
		PostExport: map[string]func(path string) ([]string, error){
			hapkg.IntegrationKind: hapkg.IntegrationPostExport,
			hapkg.PluginKind:      hapkg.PluginPostExport,
		},
	}
}
