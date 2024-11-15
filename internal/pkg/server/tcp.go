package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

func tcpServe(inChannel <-chan *http.Request, outChannel chan<- *http.Request) {
	listener, err := net.Listen("tcp", "localhost:4040")
	if err != nil {
		fmt.Println("[ERROR:TCP] ", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR:TCP:CONN] ", err)
			continue
		}

		go handleClient(conn, inChannel, outChannel)
	}
}

func handleClient(conn net.Conn, inChannel <-chan *http.Request, outChannel chan<- *http.Request) {
	defer conn.Close()

	for {
		reply := <-inChannel

		bodyBytes, err := io.ReadAll(reply.Body)
		if err != nil {
			fmt.Println("[ERROR:TCP:HTTP:REQ] ", err)
			return
		}

		_, err = conn.Write(bodyBytes)
		if err != nil {
			fmt.Println("[ERROR:TCP:CONN:WRITE] ", err)
			return
		}
	}

}
