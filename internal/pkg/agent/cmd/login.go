package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
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
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Can't open server on localhost : ", err)
		os.Exit(1)
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
	fmt.Println("You can connect to : \n\n", loginUri)

	err = srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Can't open server on localhost : ", err)
		os.Exit(1)
	}

	fmt.Println("Authorization code", authCode)
}

func createLoginUri(port string) string {
	redirectUri := "http://127.0.0.1:" + port
	encodedRedirectUri := url.QueryEscape(redirectUri)
	return BASE_URL_KEYCLOAK + "/realms/moon/protocol/openid-connect/auth?client_id=moon-agent&redirect_uri=" + encodedRedirectUri + "&response_type=code&scope=openid+email"

}
