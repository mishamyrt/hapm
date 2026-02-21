package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type listCommand struct{}

func (listCommand) New(app *hapm.App) *cobra.Command {
	listCmd := cobra.Command{
		Use:     "list",
		Short:   "List installed packages",
		Example: "hapm list",
		Args:    cobra.MaximumNArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.List()
		},
	}

	return &listCmd
}
