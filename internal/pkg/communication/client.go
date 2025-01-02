package communication

import (
	"bufio"
	"bytes"
	"crypto/tls"
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

	var bytesReceived int
	length := uint64(HEADER_SIZE) + uint64(header.LenToken) + header.LenData
	buffer := bytes.NewBuffer(nil)

	for {
		chunk := make([]byte, READ_BUFFER_SIZE)
		read, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}

		bytesReceived += read
		buffer.Write(chunk[:read])

		if buffer.Len() == int(length) {
			break
		}
	}

	packet, err := PacketFromBytes(buffer.Bytes())
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
