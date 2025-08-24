package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdRegister() *cobra.Command {
	registerCmd := cobra.Command{
		Use:   "register",
		Short: "Register to the moon server",
		Args:  cobra.NoArgs,
		Run:   handlerRegister,
	}

	return &registerCmd
}

// TODO: to be put in parameter
const BASE_API_URL = "http://m00n.fr"

func handlerRegister(cmd *cobra.Command, args []string) {
	accessToken := oidcTokenFlow(true)
	client := &http.Client{}

	req, err := http.NewRequest("POST", BASE_API_URL + "/api/users", nil)
	if err != nil {
		fmt.Println("An error occured while creating the POST request to /user", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	_, err = client.Do(req)
	if err != nil {
		fmt.Println("An error occured while sending the POST request to /user", err)
		os.Exit(1)
	}

	fmt.Println("Successfully sent request to /user")
}
