package start

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"moon/internal/pkg/agent/cmd/auth"
	"moon/internal/pkg/agent/files"
	"moon/internal/pkg/communication"
	"moon/internal/pkg/utils"
)

func connectToServer(serverAddrPort string, urlTarget *url.URL, statistics *Statistics) error {
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
	client := communication.NewClient(conn)

	interceptSignal(client)

	err = client.SendConnectionStart(tokensCached.AccessToken)
	if err != nil {
		return err
	}

	return handleRequest(client, urlTarget, statistics)
}

func handleRequest(client *communication.Client, url *url.URL, statistics *Statistics) error {
	httpClient := &http.Client{
		Timeout: time.Minute * 5, // same as google chrome
	}

	for {
		packetRequest, err := client.Read()
		if err != nil {
			return fmt.Errorf("Error while parsing packet : %v", err)
		}

		switch packetRequest.Header.Type {
		case communication.InvalidToken:
			fmt.Println("Token has expired. Please use `moon login`.")
			os.Exit(1)
			break
		case communication.HttpRequest:
			reader := bytes.NewReader(packetRequest.Payload.Data)
			reqBufio := bufio.NewReader(reader)
			req, err := http.ReadRequest(reqBufio)
			if err != nil {
				return err
			}

			req.URL.Host = url.Host
			req.URL.Scheme = url.Scheme
			req.RequestURI = ""

			reqHeaders := make(map[string]string, 0)
			for name, values := range req.Header {
				reqHeaders[name] = strings.Join(values, ", ")
			}
			body, _ := io.ReadAll(req.Body)
			// Replace a buffer to allow the body to be read again
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			defer req.Body.Close()

			// Pass request to http handler
			call := HttpMessage{
				Request: RequestMessage{
					Method:   req.Method,
					Path:     req.URL.Path,
					Headers:  reqHeaders,
					Body:     string(body),
					Datetime: fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05")),
				},
			}

			start := time.Now()

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

			reqDuration := time.Since(start)

			err = client.SendHttpResponse(buf.Bytes())
			if err != nil {
				return err
			}

			// Pass response to http handler
			respBufio := bufio.NewReader(&buf)
			respHttp, err := http.ReadResponse(respBufio, req)
			if err == nil {
				respHeaders := make(map[string]string, 0)
				for name, values := range respHttp.Header {
					respHeaders[name] = strings.Join(values, ", ")
				}
				body, _ := io.ReadAll(respHttp.Body)
				// No need to read again the body
				defer respHttp.Body.Close()

				call.Response = ResponseMessage{
					Status:   respHttp.StatusCode,
					Headers:  respHeaders,
					Body:     string(body),
					Duration: utils.FormatDuration(reqDuration),
				}
			}

			// Send to sse handler
			statistics.HttpCalls = append(statistics.HttpCalls, call)

			// No need to send to other goroutine if no EventListener
			if statistics.EventListener > 0 {
				statistics.Event <- 1
			}
			break
		default:
			// skip this packet
			continue
		}
	}
}

func getReadyForAuth() (*auth.TokenDisk, error) {
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

	var tokensCached auth.TokenDisk
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

// Heartbeat to detect lost connection
func heartbeat() {}
