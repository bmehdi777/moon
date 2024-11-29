package cmd

import (
	"github.com/bmehdi777/moon/internal/pkg/agent/cmd/login"
	"github.com/spf13/cobra"
)

func newCmdClientRoot() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "moon",
		Short: "Create a tunnel from a local port to the world",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(newCmdStart())
	rootCmd.AddCommand(login.NewCmdLogin())

	return &rootCmd
}

func ExecuteClient() error {
	return newCmdClientRoot().Execute()
}
