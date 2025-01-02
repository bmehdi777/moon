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

	fmt.Println("header bytes : ", headerBytes)
	header, err := HeaderFromBytes(headerBytes)
	if err != nil {
		return nil, err
	}

	length := uint64(HEADER_SIZE) + uint64(header.LenToken) + header.LenData
	fmt.Println("length : ", length)
	fmt.Println("lengthData2 : ", header.LenData)
	buffer := make([]byte, length)
	_, err = reader.Read(buffer)
	if err != nil {
		return nil, err
	}

	packet, err := PacketFromBytes(buffer)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Packet : %v", packet)
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
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendInvalidToken() error {
	packet := NewPacket(InvalidToken, c.AccessToken, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}
