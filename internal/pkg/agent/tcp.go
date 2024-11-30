package agent

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/bmehdi777/moon/internal/pkg/messages"
)

func ConnectToServer(serverAddrPort string, urlTarget *url.URL) error {
	// needed to be changed in prod environment
	config := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", serverAddrPort, &config)
	if err != nil {
		return errors.New("Can't connect to the server : " + err.Error())
	}
	defer conn.Close()

	return handleRequest(conn, urlTarget)
}

func handleRequest(conn *tls.Conn, url *url.URL) error {
	interceptSignal(conn)

	err := sendAuthMessage(conn)
	if err != nil {
		return err
	}
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

func sendAuthMessage(conn *tls.Conn) error {
	msg := messages.AuthRequest{
		Version: '1',
	}
	var indexBuffer bytes.Buffer
	encoder := gob.NewEncoder(&indexBuffer)
	err := encoder.Encode(&msg)
	_, err = conn.Write(indexBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func interceptSignal(conn net.Conn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		conn.Close()
		os.Exit(1)
	}()
}
