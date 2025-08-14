package auth

import "github.com/spf13/cobra"

func NewCmdLogin() *cobra.Command {
	loginCmd := cobra.Command{
		Use:   "login",
		Short: "Login to the moon server",
		Args:  cobra.NoArgs,
		Run:   handlerLogin,
	}

	return &loginCmd
}

func handlerLogin(cmd *cobra.Command, args []string) {
	//accessToken := oidcTokenFlow(false)
	oidcTokenFlow(false)

	// it should have a refresh token (offline token)
	// store it
	// use it each time we want to *start* to have an access token
	// use this access token at the begining of the tunnel
}
