package login

import (
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// to put in a parameter
const BASE_URL_KEYCLOAK = "http://localhost:8080"

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
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Can't open server : ", err)
	}
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	codeVerifier := oauth2.GenerateVerifier()
	codeChallenge := oauth2.S256ChallengeFromVerifier(codeVerifier)

	authCode, err := getAuthorizationCode(listener, port, codeChallenge)
	if err != nil {
		fmt.Println("An error occured : ", err)
	}

	_ = getToken(authCode, codeVerifier, "http://127.0.0.1:"+port)
	// now we need to store the access_token and refresh_token on disk
}
