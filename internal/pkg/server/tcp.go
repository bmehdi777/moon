package server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bmehdi777/moon/internal/pkg/server/config"
)

func tcpServe(channelsDomains ChannelsDomains) {
	cert, err := tls.LoadX509KeyPair(config.GlobalConfig.CertPemPath, config.GlobalConfig.CertKeyPath)
	if err != nil {
		log.Fatalf("TLS - %v", err)
		return
	}
	configTls := tls.Config{Certificates: []tls.Certificate{cert}}
	configTls.Rand = rand.Reader

	fullAddrFmt := fmt.Sprintf("%v:%v", config.GlobalConfig.TcpAddr, config.GlobalConfig.TcpPort)
	listener, err := tls.Listen("tcp", fullAddrFmt, &configTls)
	if err != nil {
		log.Fatalf("TLS - TCP - %v", err)
		return
	}
	defer listener.Close()
	log.Printf("TCP Server is up at %v", fullAddrFmt)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("TCP - CONN - %v", err)
			continue
		}

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(1 * time.Hour)
		}

		go handleClient(conn, channelsDomains)
	}
}

func handleClient(conn net.Conn, channelsDomains ChannelsDomains) {
	defer conn.Close()
	log.Printf("Connection started with %v", conn.RemoteAddr())
	respBytes := make([]byte, 1024)

	for {
		//reply := <-channels.RequestChannel
		// temp
		reply := &http.Request{}

		var buf bytes.Buffer
		err := reply.Write(&buf)
		if err != nil {
			log.Fatalf("TCP - HTTP - REQ - %v", err)
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("TCP - CONN - WRITE - %v", err)
			return
		}

		_, err = conn.Read(respBytes)
		if err != nil {
			log.Fatalf("TCP - CONN - READ - %v", err)
			return
		}

		reader := bytes.NewReader(respBytes)
		respBufio := bufio.NewReader(reader)
		// temp
		//resp, err := http.ReadResponse(respBufio, reply)
		_, err = http.ReadResponse(respBufio, reply)
		if err != nil {
			log.Fatalf("TCP - READER - %v", err)
			return
		}

		// temp
		//channels.ResponseChannel <- resp
	}

}
