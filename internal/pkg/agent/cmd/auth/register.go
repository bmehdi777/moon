package auth

import "github.com/spf13/cobra"

func NewCmdRegister() *cobra.Command {
	registerCmd := cobra.Command{
		Use: "register",
		Short: "Register to the moon server",
		Args: cobra.NoArgs,
		Run: handlerRegister,
	}

	return &registerCmd
}

func handlerRegister(cmd *cobra.Command, args []string) {
	oidcFlow(true)
}
