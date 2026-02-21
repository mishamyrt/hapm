package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type syncCommand struct{}

func (syncCommand) New(app *hapm.App) *cobra.Command {
	allowUnstable := false

	syncCmd := cobra.Command{
		Use:     "sync",
		Short:   "Synchronize storage with manifest",
		Example: "hapm sync",
		Args:    cobra.MaximumNArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.Sync(hapm.SyncOptions{AllowUnstable: allowUnstable})
		},
	}

	syncCmd.Flags().BoolVarP(
		&allowUnstable,
		"allow-unstable",
		"u",
		false,
		"Removes the restriction to stable versions when searching for updates",
	)

	return &syncCmd
}
