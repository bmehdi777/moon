package auth

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"golang.org/x/oauth2"
	"os"
	"strconv"
)

// to put in a parameter
const BASE_URL_KEYCLOAK = "http://localhost:8081"

func oidcTokenFlow(register bool) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Can't open server : ", err)
	}
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	codeVerifier := oauth2.GenerateVerifier()
	codeChallenge := oauth2.S256ChallengeFromVerifier(codeVerifier)

	authCode, err := getAuthorizationCode(listener, port, codeChallenge, register)
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

	// save cached tokens to disk
	//err = files.SaveToConfigFile(files.AUTH_FILENAME, diskTokenBytes)
	//if err != nil {
	//	fmt.Println("Can't save authentification data to disk : ", err)
	//}

	return keycloakJWT.AccessToken
}

func getAuthorizationCode(listener net.Listener, port string, challenge string, register bool) (string, error) {
	mux := http.NewServeMux()
	srv := http.Server{
		Handler: mux,
	}

	var authCode string

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode = r.URL.Query().Get("code")
		if authCode != "" {
			w.Write([]byte("Successfully logged in."))
		} else {
			w.Write([]byte("An error occured while trying to log you. Try later."))
		}
		go srv.Shutdown(context.Background())
	})

	uri := createAuthUri(challenge, port, register)
	fmt.Println("If your browser didn't open, you can click on the following link : \n\n", uri)

	openInBrowser(uri)

	err := srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return "", err
	}

	return authCode, nil
}

func getToken(authCode string, verifier string, callbackUri string) ([]byte, error) {
	encodedRedirectUri := url.QueryEscape(callbackUri)
	urlToken := BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/token"

	var payloadString strings.Builder
	payloadString.WriteString("grant_type=authorization_code")
	payloadString.WriteString("&client_id=moon-agent")
	payloadString.WriteString("&code=" + authCode)
	payloadString.WriteString("&code_verifier=" + verifier)
	payloadString.WriteString("&redirect_uri=" + encodedRedirectUri)

	payload := strings.NewReader(payloadString.String())

	res, err := http.Post(urlToken, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func createAuthUri(challenge string, port string, register bool) string {
	redirectUri := "http://127.0.0.1:" + port
	encodedRedirectUri := url.QueryEscape(redirectUri)

	path := "auth"
	if register {
		path = "registrations"
	}

	return BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/" + path + "?client_id=moon-agent&redirect_uri=" + encodedRedirectUri + "&response_type=code&scope=openid&code_challenge_method=S256&code_challenge=" + challenge
}
