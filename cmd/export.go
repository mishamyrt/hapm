package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type exportCommand struct{}

func (exportCommand) New(app *hapm.App) *cobra.Command {
	exportCmd := cobra.Command{
		Use:     "export <path>",
		Short:   "Export packages to Home Assistant directory structure",
		Example: "hapm export ./export",
		Args:    cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return app.Export(args)
		},
	}

	return &exportCmd
}
