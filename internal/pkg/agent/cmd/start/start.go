package start

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdStart() *cobra.Command {
	startCmd := cobra.Command{
		Use:   "start <URL>",
		Short: "Connect <URL> to the world",
		Args:  cobra.RangeArgs(1, 1),
		Run:   handlerStart,
	}

	startCmd.PersistentFlags().String("server-address", "", "Specify the server address to connect")
	startCmd.PersistentFlags().String("server-port", "", "Specify the server port to connect")

	return &startCmd
}

func handlerStart(cmd *cobra.Command, args []string) {
	urlTarget, err := url.ParseRequestURI(args[0])
	if err != nil {
		fmt.Println("The URL provided isn't valid.")
		os.Exit(1)
	}

	addr, _ := cmd.PersistentFlags().GetString("server-address")
	port, _ := cmd.PersistentFlags().GetString("server-port")

	// debug purpose
	if addr == "" && port == "" {
		addr = "localhost"
		port = "4040"
	}

	stats := make(Statistics)

	go httpServe(&stats)

	err = connectToServer(addr+":"+port, urlTarget, &stats)
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
}

