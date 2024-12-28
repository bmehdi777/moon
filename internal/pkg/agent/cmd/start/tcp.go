package start

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/bmehdi777/moon/internal/pkg/agent/cmd/login"
	"github.com/bmehdi777/moon/internal/pkg/agent/files"
	"github.com/bmehdi777/moon/internal/pkg/messages"
)

func connectToServer(serverAddrPort string, urlTarget *url.URL) error {
	tokensCached, err := getReadyForAuth()
	if err != nil {
		return err
	}
	// needed to be changed in prod environment
	config := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", serverAddrPort, &config)
	if err != nil {
		return errors.New("Can't connect to the server : " + err.Error())
	}
	defer conn.Close()

	// TODO: fix ctrl-c doesnt close connection (precisely, it close but still use it afterwards)
	interceptSignal(conn)

	err = sendAuth(conn, *tokensCached)
	if err != nil {
		return err
	}

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

func getReadyForAuth() (*login.TokenDisk, error) {
	/*
	* Get tokens from cache :
	* If no auth file exist, user should login
	* If expires are outdated, remove the cache and user should login
	* If ANY error, remove the cache file and user should login
	 */

	if exist := files.IsFromConfigFileExist(files.AUTH_FILENAME); !exist {
		return nil, errors.New("No session exist.\n\nUse 'moon login' before starting a tunnel.")
	}

	tokensCachedBytes, err := files.ReadFromConfigFile(files.AUTH_FILENAME)
	if err != nil {
		return nil, err
	}

	var tokensCached login.TokenDisk
	err = json.Unmarshal(tokensCachedBytes, &tokensCached)
	if err != nil {
		return nil, errors.New("Can't parse auth file. It should be a JSON file.")
	}

	n := time.Now()
	atExpTimestamp := time.Unix(tokensCached.AccessTokenExpire, 0)
	rtExpTimestamp := time.Unix(tokensCached.RefreshTokenExpire, 0)

	if n.After(atExpTimestamp) && n.After(rtExpTimestamp) {
		return nil, errors.New("Your session expired. Use 'moon login' to start a new session.")
	} else if n.After(atExpTimestamp) && n.Before(rtExpTimestamp) {
		// TODO: here we should ask another access token with the refresh token
	}
	return &tokensCached, nil
}

func sendAuth(conn *tls.Conn, tokensCached login.TokenDisk) error {

	// send the message to the server
	msg := messages.AuthRequest{
		Version:        '1',
		AccessTokenJWT: tokensCached.AccessToken,
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