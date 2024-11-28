package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
)

const BASE_URL_KEYCLOAK = "http://localhost:8080"

// http://localhost:8080/realms/moon/protocol/openid-connect/auth?client_id=moon-test&redirect_uri=https%3A%2F%2Fwww.keycloak.org%2Fapp%2F%23url%3Dhttp%3A%2F%2Flocalhost%3A8080%26realm%3Dmoon%26client%3Dmoon-test&state=0b23d3a6-b15f-4e6c-bd29-06b96bed1fd3&response_mode=fragment&response_type=code&scope=openid
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
	authCode, err := getAuthorizationCode()
	if err != nil {
		fmt.Println("An error occured : ", err)
	}

	fmt.Println("Authorization code : ", authCode)
}

func getAuthorizationCode() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}

	mux := http.NewServeMux()
	srv := http.Server{
		Handler: mux,
	}

	var authCode string

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Successfully logged in."))
		authCode = r.URL.Query().Get("code")
		go srv.Shutdown(context.Background())
	})

	port := listener.Addr().(*net.TCPAddr).Port
	loginUri := createLoginUri(strconv.Itoa(port))
	fmt.Println("If your browser isn't open, you can click on the following link : \n\n", loginUri)

	openInBrowser(loginUri)

	err = srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return "", err
	}

	return authCode, nil
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

func createLoginUri(port string) string {
	redirectUri := "http://127.0.0.1:" + port
	encodedRedirectUri := url.QueryEscape(redirectUri)
	return BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/auth?client_id=moon-agent&redirect_uri=" + encodedRedirectUri + "&response_type=code&scope=openid+email"

}
