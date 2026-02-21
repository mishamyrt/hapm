package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type updatesCommand struct{}

func (updatesCommand) New(app *hapm.App) *cobra.Command {
	allowUnstable := false

	updatesCmd := cobra.Command{
		Use:     "updates",
		Short:   "Show available package updates",
		Example: "hapm updates",
		Args:    cobra.MaximumNArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.PrintUpdates(hapm.UpdatesOptions{AllowUnstable: allowUnstable})
		},
	}

	updatesCmd.Flags().BoolVarP(
		&allowUnstable,
		"allow-unstable",
		"u",
		false,
		"Removes the restriction to stable versions when searching for updates",
	)

	return &updatesCmd
}
