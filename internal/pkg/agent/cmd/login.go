package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

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

	fmt.Println("Authorization code : ", authCode)
	bodyResponse := getToken(authCode, codeVerifier, "http://127.0.0.1:"+port)

	fmt.Println("\nBody response : ", bodyResponse)
}

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
	fmt.Println("If your browser isn't open, you can click on the following link : \n\n", loginUri)

	openInBrowser(loginUri)

	err := srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return "", err
	}

	return authCode, nil
}

func getToken(authCode string, verifier string, callbackUri string) string {
	encodedRedirectUri := url.QueryEscape(callbackUri)
	url := BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/token"

	var payloadString strings.Builder
	payloadString.WriteString("grant_type=authorization_code")
	payloadString.WriteString("&client_id=moon-agent")
	payloadString.WriteString("&code=" + authCode)
	payloadString.WriteString("&code_verifier=" + verifier)
	payloadString.WriteString("&redirect_uri=" + encodedRedirectUri)

	fmt.Println("DEBUG : ", payloadString.String())
	payload := strings.NewReader(payloadString.String())

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}

func createLoginUri(challenge string, port string) string {
	redirectUri := "http://127.0.0.1:" + port
	encodedRedirectUri := url.QueryEscape(redirectUri)

	return BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/auth?client_id=moon-agent&redirect_uri=" + encodedRedirectUri + "&response_type=code&scope=openid+email&code_challenge_method=S256&code_challenge=" + challenge
}

func openInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}

	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return err
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unsupported platform")
	}
}
