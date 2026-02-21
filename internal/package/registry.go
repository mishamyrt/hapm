package hapmpkg

import "io"

type Constructor func(description PackageDescription, rootPath string, client GitClient) Package

type Registry struct {
	Constructors map[string]Constructor
	PreExport    map[string]func(path string) error
	PostExport   map[string]func(path string, out io.Writer) error
}

func DefaultRegistry() Registry {
	return Registry{
		Constructors: map[string]Constructor{
			IntegrationKind: NewIntegrationPackage,
			PluginKind:      NewPluginPackage,
		},
		PreExport: map[string]func(path string) error{
			IntegrationKind: IntegrationPreExport,
			PluginKind:      PluginPreExport,
		},
		PostExport: map[string]func(path string, out io.Writer) error{
			IntegrationKind: IntegrationPostExport,
			PluginKind:      PluginPostExport,
		},
	}
}

func (r Registry) SupportedKinds() []string {
	kinds := make([]string, 0, len(r.Constructors))
	for kind := range r.Constructors {
		kinds = append(kinds, kind)
	}
	return kinds
}
