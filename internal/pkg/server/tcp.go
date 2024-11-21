package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"
)

func tcpServe(channelsDomains ChannelsDomains) {
	listener, err := net.Listen("tcp", ":4040")
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

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(1 * time.Hour)
		}

		go handleClient(conn, channelsDomains)
	}
}

func negotiateConnection(conn net.Conn) {

}

func handleClient(conn net.Conn, channelsDomains ChannelsDomains) {
	defer conn.Close()
	fmt.Println("Connection started")
	respBytes := make([]byte, 1024)

	for {
		//reply := <-channels.RequestChannel
		// temp
		reply := &http.Request{}

		var buf bytes.Buffer
		err := reply.Write(&buf)
		if err != nil {
			fmt.Println("[ERROR:TCP:HTTP:REQ] ", err)
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			fmt.Println("[ERROR:TCP:CONN:WRITE] ", err)
			return
		}

		_, err = conn.Read(respBytes)
		if err != nil {
			fmt.Println("[ERROR:TCP:CONN:READ]", err)
			return
		}

		reader := bytes.NewReader(respBytes)
		respBufio := bufio.NewReader(reader)
		// temp
		//resp, err := http.ReadResponse(respBufio, reply)
		_, err = http.ReadResponse(respBufio, reply)
		if err != nil {
			fmt.Println("[ERROR:TCP:READER]", err)
			return
		}

		// temp
		//channels.ResponseChannel <- resp
	}

}
