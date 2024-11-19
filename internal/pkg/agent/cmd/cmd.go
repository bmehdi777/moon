package cmd

import "github.com/spf13/cobra"

func newCmdClientRoot() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "moon",
		Short: "Create a tunnel from a local port to the world",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(newCmdStart())

	return &rootCmd
}

func ExecuteClient() error {
	return newCmdClientRoot().Execute()
}
