package communication

import (
	"crypto/tls"
)

type Client struct {
	Connection  *tls.Conn
	AccessToken string
}

func NewClient(conn *tls.Conn, token string) *Client {
	return &Client{Connection: conn, AccessToken: token}
}

func (c *Client) sendPacket(packet *Packet) error {
	_, err := c.Connection.Write(packet.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendConnectionStart() error {
	packet := NewPacket(ConnectionStart, c.AccessToken, nil)
	err := c.sendPacket(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendConnectionClose() error {
	packet := NewPacket(ConnectionClose, c.AccessToken, nil)
	err := c.sendPacket(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHttpRequest(data []byte) error {
	packet := NewPacket(HttpRequest, c.AccessToken, data)
	err := c.sendPacket(packet)
	if err != nil {
		return err
	}
	return nil
}
