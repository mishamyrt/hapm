package cmd

import (
	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type versionsCommand struct{}

func (versionsCommand) New(app *hapm.App) *cobra.Command {
	versionsCmd := cobra.Command{
		Use:     "versions <location>",
		Short:   "List available versions for a package",
		Example: "hapm versions foo/bar",
		Args:    cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return app.PrintVersions(args)
		},
	}

	return &versionsCmd
}
