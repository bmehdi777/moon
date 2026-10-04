package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"moon/internal/pkg/agent/files"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

// to put in a parameter
const BASE_URL_KEYCLOAK = "http://localhost:8081"

func oidcTokenFlow(register bool) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Can't open server : ", err)
		return ""
	}
	defer listener.Close()
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		fmt.Println("Can't determine callback port")
		return ""
	}
	port := strconv.Itoa(tcpAddr.Port)

	codeVerifier := oauth2.GenerateVerifier()
	codeChallenge := oauth2.S256ChallengeFromVerifier(codeVerifier)

	authCode, err := getAuthorizationCode(listener, port, codeChallenge, register)
	if err != nil {
		fmt.Println("An error occured : ", err)
		return ""
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

	// todo: encode this with a local password ??
	err = files.SaveToConfigFile(files.AUTH_FILENAME, diskTokenBytes)
	if err != nil {
		fmt.Println("Can't save authentification data to disk : ", err)
	}

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
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("Keycloak token request failed with status %s: %w", res.Status, readErr)
		}
		return nil, fmt.Errorf("Keycloak token request failed with status %s: %s", res.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func doRefreshToken(refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is empty")
	}
	endpoint := oauth2.Endpoint{TokenURL: BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/token"}
	config := oauth2.Config{ClientID: "moon-agent", Endpoint: endpoint}
	token, err := config.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken}).Token()
	if err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("Keycloak returned an empty access token")
	}
	if token.RefreshToken == "" {
		token.RefreshToken = refreshToken
	}
	return token, nil
}

// RefreshToken exchanges a cached refresh token for a new access token.
func RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return doRefreshToken(refreshToken)
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
