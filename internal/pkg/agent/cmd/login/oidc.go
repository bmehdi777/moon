package login

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)


func getAuthorizationCode(listener net.Listener, port string, challenge string) (string, error) {
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

	loginUri := createLoginUri(challenge, port)
	fmt.Println("If your browser didn't open, you can click on the following link : \n\n", loginUri)

	openInBrowser(loginUri)

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

	res, err := http.Post( urlToken,"application/x-www-form-urlencoded", payload)
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

func createLoginUri(challenge string, port string) string {
	redirectUri := "http://127.0.0.1:" + port
	encodedRedirectUri := url.QueryEscape(redirectUri)

	return BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/auth?client_id=moon-agent&redirect_uri=" + encodedRedirectUri + "&response_type=code&scope=openid&code_challenge_method=S256&code_challenge=" + challenge
}

