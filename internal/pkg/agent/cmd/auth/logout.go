package auth

import "github.com/spf13/cobra"

func NewCmdLogout() *cobra.Command {
	logoutCmd := cobra.Command{
		Use: "logout",
		Short: "Logout of the moon server",
		Args: cobra.NoArgs,
		Run: handlerLogout,
	}

	return &logoutCmd
}

func handlerLogout(cmd *cobra.Command, args []string) {
	// verify the authent token is on the machine
	// if not, print err
	// if it is, use it to revoke on the server
	// then remove it from the machine 
}
