package communication

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"time"
)

type WatchdogAction string

const (
	WD_RESET WatchdogAction = "reset"
)
const WATCHDOG_TIME = 10 * time.Second

type Watchdog struct {
	Timer  *time.Timer
	Action chan WatchdogAction
}

type Client struct {
	Connection       *tls.Conn
	Watchdog         *Watchdog
}

func NewClient(conn *tls.Conn) *Client {
	return &Client{Connection: conn, Watchdog: &Watchdog{}}
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
	length := uint64(HEADER_SIZE) + header.LenData
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

func (c *Client) SendConnectionStart(token string) error {
	authMsg := NewAuthMessage(token)
	packet := NewPacket(ConnectionStart, authMsg.Bytes())
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendConnectionClose() error {
	packet := NewPacket(ConnectionClose, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendPing() error {
	packet := NewPacket(Ping, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendPong() error {
	packet := NewPacket(Pong, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHttpRequest(data []byte) error {
	packet := NewPacket(HttpRequest, data)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHttpResponse(data []byte) error {
	packet := NewPacket(HttpResponse, data)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendUnauthorized() error {
	packet := NewPacket(Unauthorized, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendAuthorized() error {
	packet := NewPacket(Authorized, nil)
	err := c.Write(packet)
	if err != nil {
		return err
	}
	return nil
}
