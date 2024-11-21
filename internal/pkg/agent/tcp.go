package agent

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
)

func ConnectToServer(serverAddrPort string, urlTarget *url.URL) error {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		return errors.New("Can't find certificates : " + err.Error())
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", serverAddrPort, &config)
	if err != nil {
		return errors.New("Can't connect to the server : " + err.Error())
	}
	defer conn.Close()

	return handleRequest(conn, urlTarget)
}

func handleRequest(conn *tls.Conn, url *url.URL) error {
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
