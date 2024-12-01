package cmd

import (
	"fmt"
	"os"

	"github.com/bmehdi777/moon/internal/pkg/agent/cmd/login"
	"github.com/bmehdi777/moon/internal/pkg/agent/cmd/start"
	"github.com/bmehdi777/moon/internal/pkg/agent/files"
	"github.com/spf13/cobra"
)

func newCmdClientRoot() *cobra.Command {
	err := files.InitConfigFolders()
	if err != nil {
		fmt.Println("An error occured : ", err)
		os.Exit(1)
	}

	rootCmd := cobra.Command{
		Use:   "moon",
		Short: "Create a tunnel from a local port to the world",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(start.NewCmdStart())
	rootCmd.AddCommand(login.NewCmdLogin())

	return &rootCmd
}

func ExecuteClient() error {
	return newCmdClientRoot().Execute()
}
