package start

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	pb "github.com/bmehdi777/moon/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	pb.TunnelClient

	ProxyUrl *url.URL
}

func (c *Client) Run(ctx context.Context) error {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c.TunnelClient = pb.NewTunnelClient(conn)

	client, err := c.TunnelClient.Stream(ctx)
	if err != nil {
		return nil
	}

	tokensCached, err := getReadyForAuth()
	if err != nil {
		return err
	}
	fmt.Println("token : ", tokensCached.AccessToken)

	err = client.Send(&pb.StreamClient{
		Event: &pb.StreamClient_Credentials{
			Credentials: &pb.StreamClient_Login{
				AccessToken: tokensCached.AccessToken,
			},
		},
	})
	if err != nil {
		fmt.Println("Error login", err)
		os.Exit(1)
	}

	fmt.Println("login")

	// heartbeat
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			fmt.Println("tick")
			select {
			case <-ticker.C:
				err := client.Send(&pb.StreamClient{
					Event: &pb.StreamClient_Ping_{},
				})
				if err != nil {
					fmt.Println("Error while ping")
					return
				}
			}
		}
	}()

	httpClient := &http.Client{
		Timeout: time.Minute * 5, // same as google chrome
	}

	for {
		msg, err := client.Recv()
		if err == io.EOF {
			fmt.Println("End of stream")
			return err
		}
		if err != nil {
			fmt.Println("Err : ", err)
			return err
		}

		switch msg.Event.(type) {
		case *pb.StreamServer_Pong_:
			fmt.Println("Pong")
			break
		case *pb.StreamServer_HttpRequest:
			reader := bytes.NewReader(msg.Event.(*pb.StreamServer_HttpRequest).HttpRequest)
			msgBufio := bufio.NewReader(reader)
			req, err := http.ReadRequest(msgBufio)
			if err != nil {
				return err
			}

			req.URL.Host = c.ProxyUrl.Host
			req.URL.Scheme = c.ProxyUrl.Scheme
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

			err = client.Send(&pb.StreamClient{
				Event: &pb.StreamClient_HttpResponse{
					HttpResponse: buf.Bytes(),
				},
			})
			if err != nil {
				return err
			}
			break
		}

	}

}
