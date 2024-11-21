package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmehdi777/moon/internal/pkg/server/database"
)

type ChannelsHttp struct {
	RequestChannel  chan *http.Request
	ResponseChannel chan *http.Response
}

type ChannelsDomains map[string]ChannelsHttp

func Run() {
	db, err := database.InitializeDBConn()
	if err != nil {
		fmt.Println("Can't connect to the database.")
		os.Exit(1)
	}

	channelPerDomain := make(ChannelsDomains)

	// channel should be created when tcp connection is created
	tcpHttpChannel := make(chan *http.Response)
	httpTcpChannel := make(chan *http.Request)

	// tcp connection between client and server
	go tcpServe(channelPerDomain)

	// http connection between server and the world
	if err := httpServe(tcpHttpChannel, httpTcpChannel, db); err != nil {
		fmt.Println("[ERROR:HTTP] ", err)
	}
}
