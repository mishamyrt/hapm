package cmd

import (
	"errors"
	"io"

	"github.com/mishamyrt/hapm/internal/hapm"
	"github.com/spf13/cobra"
)

type command interface {
	New(*hapm.App) *cobra.Command
}

var commands = [...]command{
	initCommand{},
	syncCommand{},
	installCommand{},
	updatesCommand{},
	versionsCommand{},
	listCommand{},
	exportCommand{},
}

func newRootCmd(stdout io.Writer, stderr io.Writer) *cobra.Command {
	app := hapm.New(stdout, stderr)
	globals := hapm.DefaultGlobalOptions()

	rootCmd := &cobra.Command{
		Use:   "hapm",
		Short: "Home Assistant Package Manager",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			app.SetGlobals(globals)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			_ = cmd.Help()
			return hapm.HandledError(errors.New("command is required"))
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&globals.Manifest, "manifest", "m", globals.Manifest, "Manifest path")
	rootCmd.PersistentFlags().StringVarP(&globals.Storage, "storage", "s", globals.Storage, "Storage location")
	rootCmd.PersistentFlags().BoolVarP(
		&globals.Dry,
		"dry",
		"d",
		globals.Dry,
		"Only print information. Do not make any changes to the files",
	)

	for _, command := range commands {
		rootCmd.AddCommand(command.New(app))
	}

	return rootCmd
}
