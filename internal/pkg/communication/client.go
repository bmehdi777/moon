package communication

import (
	"crypto/tls"
)

type Client struct {
	Connection *tls.Conn
}

func NewClient(conn *tls.Conn) *Client {
	return &Client{Connection: conn}
}

func (c *Client) SendConnectionStart() {}
func (c *Client) SendConnectionClose() {}
func (c *Client) SendHttpRequest() {}
