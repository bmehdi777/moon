package communication

import (
	"bufio"
	"crypto/tls"
	"fmt"
)

type Client struct {
	Connection  *tls.Conn
	AccessToken *string
}

func NewClient(conn *tls.Conn, token *string) *Client {
	return &Client{Connection: conn, AccessToken: token}
}

func (c *Client) Read() (*Packet, error) {
	// determine header size
	reader := bufio.NewReader(c.Connection)
	headerBytes, err := reader.Peek(int(HEADER_SIZE))
	if err != nil {
		return nil, err
	}

	header, err := HeaderFromBytes(headerBytes)
	if err != nil {
		return nil, err
	}

	length := HEADER_SIZE + header.LenToken + header.LenData
	buffer := make([]byte, length)
	_, err = reader.Read(buffer)
	if err != nil {
		return nil, err
	}

	packet, err := PacketFromBytes(buffer)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func (c *Client) Write(packet *Packet) error {
	_, err := c.Connection.Write(packet.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendConnectionStart() error {
	packet := NewPacket(ConnectionStart, c.AccessToken, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendConnectionClose() error {
	fmt.Println("Closing connection")
	packet := NewPacket(ConnectionClose, c.AccessToken, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHttpRequest(data []byte) error {
	packet := NewPacket(HttpRequest, c.AccessToken, data)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHttpResponse(data []byte) error {
	packet := NewPacket(HttpResponse, c.AccessToken, data)
	fmt.Printf("packet : %#v", packet)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}
