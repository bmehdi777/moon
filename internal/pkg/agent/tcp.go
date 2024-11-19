package agent

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
	"net/url"
)

func ConnectToServer(urlTarget *url.URL) error {
	serverAddr := "localhost:4040"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return errors.New("Can't reach the server. Maybe it's down ?")
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return errors.New("Can't connect to the server.")
	}
	defer conn.Close()

	return handleRequest(conn, urlTarget)
}

func handleRequest(conn *net.TCPConn, url *url.URL) error {
	httpClient := &http.Client{}
	for {
		msgBytes := make([]byte, 1024)
		_, err := conn.Read(msgBytes)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(msgBytes)
		msgBufio := bufio.NewReader(reader)
		req, err := http.ReadRequest(msgBufio)
		if err != nil {
			return err
		}

		req.URL = url
		req.RequestURI = ""

		// send to urlTarget
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = resp.Write(&buf)
		if err != nil {
			return err
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}
}
