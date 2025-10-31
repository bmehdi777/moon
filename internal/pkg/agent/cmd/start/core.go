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
	"moon/internal/pkg/server/config"
	"moon/internal/pkg/utils"
)

func connectToServer(serverAddrPort string, urlTarget *url.URL, statistics *Statistics) error {
	jwt, err := prepareAuth()
	if err != nil {
		return err
	}

	// TODO: needed to be changed in prod environment
	var tlsConfig tls.Config
	if config.GlobalConfig.Env == config.ENV_PROD {
		tlsConfig = tls.Config{InsecureSkipVerify: false}
	} else {
		tlsConfig = tls.Config{InsecureSkipVerify: true}
	}

	conn, err := tls.Dial("tcp", serverAddrPort, &tlsConfig)
	if err != nil {
		return errors.New("Can't connect to the server : " + err.Error())
	}
	defer conn.Close()
	client := communication.NewClient(conn)

	interceptSignal(client)

	err = client.SendConnectionStart(jwt.AccessToken)
	if err != nil {
		return err
	}

	return handleConnection(client, urlTarget, statistics)
}

func handleConnection(client *communication.Client, url *url.URL, statistics *Statistics) error {
	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}

	heartbeat(client)

	for {
		packetRequest, err := client.Read()
		if err != nil {
			return fmt.Errorf("Error while parsing packet : %v", err)
		}

		switch packetRequest.Header.Type {
		case communication.Unauthorized:
			fmt.Println("Token has expired. Please use `moon login`.")
			os.Exit(1)
			break
		case communication.Authorized:
			fmt.Println("Successfully connected to the server.")
			break
		case communication.HttpRequest:

			messageLocalApi := &HttpMessage{}

			req, resp, reqDuration, err := sendRequestToTarget(packetRequest, url, httpClient, messageLocalApi)

			var buf bytes.Buffer
			err = resp.Write(&buf)
			if err != nil {
				return err
			}

			// send resp back to moon server (internet)
			err = client.SendHttpResponse(buf.Bytes())
			if err != nil {
				return err
			}

			sendStatsToLocal(&buf, req, messageLocalApi, reqDuration, statistics)

			break
		case communication.Pong:
		default:
			// skip this packet
			continue
		}
	}
}

func sendRequestToTarget(packetRequest *communication.Packet, url *url.URL, httpClient *http.Client, message *HttpMessage) (*http.Request, *http.Response, *time.Duration, error) {
	// parse binary packet to http request
	reader := bytes.NewReader(packetRequest.Payload.Data)
	reqBufio := bufio.NewReader(reader)
	req, err := http.ReadRequest(reqBufio)
	if err != nil {
		return nil, nil, nil, err
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
	message.Request = RequestMessage{
		Method:   req.Method,
		Path:     req.URL.Path,
		Headers:  reqHeaders,
		Body:     string(body),
		Datetime: fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05")),
	}

	start := time.Now()

	// send to url target
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}

	requestDuration := time.Since(start)

	return req, resp, &requestDuration, nil
}
func sendStatsToLocal(respBytes *bytes.Buffer, req *http.Request, message *HttpMessage, requestDuration *time.Duration, statistics *Statistics) {
	// Pass response to local api
	respBufio := bufio.NewReader(respBytes)
	respHttp, err := http.ReadResponse(respBufio, req)
	if err == nil {
		respHeaders := make(map[string]string, 0)
		for name, values := range respHttp.Header {
			respHeaders[name] = strings.Join(values, ", ")
		}
		body, _ := io.ReadAll(respHttp.Body)
		// No need to read again the body
		defer respHttp.Body.Close()

		message.Response = ResponseMessage{
			Status:   respHttp.StatusCode,
			Headers:  respHeaders,
			Body:     string(body),
			Duration: utils.FormatDuration(*requestDuration),
		}
	}

	// Send to SSE handler
	statistics.HttpCalls = append(statistics.HttpCalls, *message)

	// No need to send to other goroutine if no EventListener
	if statistics.EventListener > 0 {
		statistics.Event <- 1
	}

}

func prepareAuth() (*auth.TokenDisk, error) {
	/*
	* Get refresh token from cache :
	* If no auth file exist, user should login
	* If expires are outdated, remove the cache and user should login
	* If ANY error, remove the cache file and user should login
	 */

	if exist := files.IsFromConfigFileExist(files.AUTH_FILENAME); !exist {
		return nil, errors.New("No session exist.\n\nUse 'moon login' before starting a tunnel.")
	}

	refreshTokensCachedBytes, err := files.ReadFromConfigFile(files.AUTH_FILENAME)
	if err != nil {
		return nil, err
	}

	var tokensCached auth.TokenDisk
	err = json.Unmarshal(refreshTokensCachedBytes, &tokensCached)
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
func heartbeat(client *communication.Client) {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := client.SendPing()
				if err != nil {
					fmt.Println("Can't contact server")
					os.Exit(1)
				}
			}
		}
	}()
}
