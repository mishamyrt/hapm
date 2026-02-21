package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type installCommand struct{}

func (installCommand) New(app *hapm.App) *cobra.Command {
	allowUnstable := false
	packageType := ""

	installCmd := cobra.Command{
		Use:     "install",
		Short:   "Install new packages or update existing ones",
		Example: "hapm install --type integrations foo/bar@v1.0.0",
		Args:    cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return app.Install(hapm.InstallOptions{
				Entries:       args,
				PackageType:   packageType,
				AllowUnstable: allowUnstable,
			})
		},
	}

	installCmd.Flags().StringVarP(
		&packageType,
		"type",
		"t",
		"",
		"Packages type. Required parameter if a new package is installed",
	)
	installCmd.Flags().BoolVarP(
		&allowUnstable,
		"allow-unstable",
		"u",
		false,
		"Removes the restriction to stable versions when searching for updates",
	)

	return &installCmd
}
