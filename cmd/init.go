package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type initCommand struct{}

func (initCommand) New(app *hapm.App) *cobra.Command {
	initCmd := cobra.Command{
		Use:     "init",
		Short:   "Initialize an empty manifest",
		Example: "hapm init",
		Args:    cobra.MaximumNArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.Init()
		},
	}

	return &initCmd
}
