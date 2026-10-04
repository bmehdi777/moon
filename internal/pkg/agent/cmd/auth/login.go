package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdLogin() *cobra.Command {
	loginCmd := cobra.Command{
		Use:   "login",
		Short: "Login to the moon server",
		Args:  cobra.NoArgs,
		Run:   handlerLogin,
	}
	loginCmd.Flags().String("auth-server", BASE_URL_KEYCLOAK, "Keycloak server URL")

	return &loginCmd
}

func handlerLogin(cmd *cobra.Command, args []string) {
	authServer, err := cmd.Flags().GetString("auth-server")
	if err != nil {
		fmt.Println("Can't read auth server: ", err)
		return
	}
	//accessToken := oidcTokenFlow(false)
	oidcTokenFlow(authServer, false)

	// it should have a refresh token (offline token)
	// store it
	// use it each time we want to *start* to have an access token
	// use this access token at the begining of the tunnel
}
