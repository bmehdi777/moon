package login

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"moon/internal/pkg/agent/files"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// to put in a parameter
const BASE_URL_KEYCLOAK = "http://localhost:8081"

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

	tokenResponse, err := getToken(authCode, codeVerifier, "http://127.0.0.1:"+port)
	if err != nil {
		fmt.Println("An error occured while getting token : ", err)
		os.Exit(1)
	}

	var keycloakJWT KeycloakJWTS
	err = json.Unmarshal(tokenResponse, &keycloakJWT)
	if err != nil {
		fmt.Println("An error occured while parsing the token ", err)
		os.Exit(1)
	}

	diskTokenBytes, err := json.Marshal(keycloakJWT.ToDisk())
	if err != nil {
		fmt.Println("Can't parse to json disk token : ", err)
		os.Exit(1)
	}

	// save cached tokens to disk
	err = files.SaveToConfigFile(files.AUTH_FILENAME, diskTokenBytes)
	if err != nil {
		fmt.Println("Can't save authentification data to disk : ", err)
	}
}
