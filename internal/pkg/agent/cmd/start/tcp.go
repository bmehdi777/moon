package start

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/bmehdi777/moon/internal/pkg/agent/cmd/login"
	"github.com/bmehdi777/moon/internal/pkg/agent/files"
	"github.com/bmehdi777/moon/internal/pkg/communication"
)

func connectToServer(serverAddrPort string, urlTarget *url.URL, stats *Statistics) error {
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
	client := communication.NewClient(conn, &tokensCached.AccessToken)

	interceptSignal(client)

	err = client.SendConnectionStart()
	if err != nil {
		return err
	}

	return handleRequest(client, urlTarget, stats)
}

func handleRequest(client *communication.Client, url *url.URL, stats *Statistics) error {
	httpClient := &http.Client{
		Timeout: time.Minute * 5, // same as google chrome
	}

	for {
		packetRequest, err := client.Read()
		if err != nil {
			return fmt.Errorf("Error while parsing packet : %v", err)
		}

		if packetRequest.Header.Type == communication.InvalidToken {
			fmt.Println("Token has expired. Please use `moon login`.")
			os.Exit(1)
		}

		if packetRequest.Header.Type != communication.HttpRequest {
			// skip this packet if it isn't a request
			continue
		}

		reader := bytes.NewReader(packetRequest.Payload.Data)
		msgBufio := bufio.NewReader(reader)
		req, err := http.ReadRequest(msgBufio)
		if err != nil {
			return err
		}

		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.RequestURI = ""

		statsKey := req.Method + " " + req.URL.Path
		if stats.Get(statsKey) == -1 {
			stats.Set(statsKey, 1)
		} else {
			stats.Increment(statsKey)
		}

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

		err = client.SendHttpResponse(buf.Bytes())
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

func interceptSignal(client *communication.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		client.SendConnectionClose()
		os.Exit(1)
	}()
}
