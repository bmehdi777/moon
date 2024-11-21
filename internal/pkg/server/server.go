package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/bmehdi777/moon/internal/pkg/server/database"
)

type ChannelsHttp struct {
	RequestChannel  chan *http.Request
	ResponseChannel chan *http.Response
}

type ChannelsDomains map[string]ChannelsHttp

func Run() {
	config.InitConfig()

	db, err := database.InitializeDBConn()
	if err != nil {
		fmt.Println("Can't connect to the database.")
		os.Exit(1)
	}

	channelPerDomain := make(ChannelsDomains)

	// tcp connection between client and server
	go tcpServe(channelPerDomain)

	// http connection between server and the world
	if err := httpServe(channelPerDomain, db); err != nil {
		fmt.Println("[ERROR:HTTP] ", err)
	}
}
