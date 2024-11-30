package login

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
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
	//fmt.Println("Token : ", tokenResponse)
	// now we need to store the access_token and refresh_token on disk
	var t KeycloakToken
	err = json.Unmarshal(tokenResponse, &t)
	if err != nil {
		fmt.Println("An error occured while parsing the token ", err)
		os.Exit(1)
	}

}

type KeycloakToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

