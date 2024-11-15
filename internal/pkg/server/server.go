package server

import (
	"fmt"
	"net/http"
)

func Run() {

	// channel should be created when tcp connection is created
	tcpHttpChannel := make(chan *http.Request)
	httpTcpChannel := make(chan *http.Request)

	// tcp connection between client and server
	go tcpServe(httpTcpChannel, tcpHttpChannel)

	// http connection between server and the world
	if err := httpServe(tcpHttpChannel, httpTcpChannel); err != nil {
		fmt.Println("[ERROR:HTTP] ", err)
	}
}
