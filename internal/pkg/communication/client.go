package communication

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
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
	Connection *tls.Conn
	Watchdog   *Watchdog
	reader     *bufio.Reader
}

func NewClient(conn *tls.Conn) *Client {
	return &Client{Connection: conn, Watchdog: &Watchdog{}, reader: bufio.NewReader(conn)}
}

func (c *Client) Read() (*Packet, error) {
	headerBytes := make([]byte, HEADER_SIZE)
	if _, err := io.ReadFull(c.reader, headerBytes); err != nil {
		return nil, err
	}
	header, err := HeaderFromBytes(headerBytes)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(headerBytes)
	if header.LenData > 0 {
		payload := make([]byte, int(header.LenData))
		if _, err := io.ReadFull(c.reader, payload); err != nil {
			return nil, err
		}
		buffer.Write(payload)
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

func (c *Client) SendHttpRequest(url string, data []byte) error {
	httpReqMsg := NewHttpRequestMessage(url, data)
	packet := NewPacket(HttpRequest, httpReqMsg.Bytes())
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
