package communication

import (
	"encoding/binary"
	"fmt"
)

// Message type is a field in the header. It allows us to conver the payload
// to the right struct.
type MessageType uint8

func (mt MessageType) String() string {
	switch mt {
	case ConnectionStart:
		return "ConnectionStart"
	case ConnectionClose:
		return "ConnectionClose"
	case Ping:
		return "Ping"
	case Pong:
		return "Pong"
	case HttpRequest:
		return "HttpRequest"
	case HttpResponse:
		return "HttpResponse"
	case Unauthorized:
		return "Unauthorized"
	case Authorized:
		return "Authorized"
	default:
		return fmt.Sprintf("%d", mt)
	}
}

const (
	ConnectionStart MessageType = iota
	ConnectionClose
	// Heartbeat
	Ping
	Pong

	HttpRequest
	HttpResponse

	Authorized
	Unauthorized
)

// Complexe payload are converted to type
type Message interface {
	Bytes() []byte
}

// The server send the public url to the client
// The client will be able to swap location field if needed
type HttpRequestMessage struct {
	UrlLength  uint32
	Url        string
	DataLength uint32
	Data       []byte
}

func NewHttpRequestMessage(url string, data []byte) HttpRequestMessage {
	return HttpRequestMessage{
		UrlLength:  uint32(len(url)),
		Url:        url,
		DataLength: uint32(len(data)),
		Data:       data,
	}
}

func (hrm *HttpRequestMessage) Bytes() []byte {
	msg := make([]byte, 0)

	msg = binary.BigEndian.AppendUint32(msg, hrm.UrlLength)

	msg = append(msg, []byte(hrm.Url)...)

	msg = binary.BigEndian.AppendUint32(msg, hrm.DataLength)

	msg = append(msg, hrm.Data...)

	return msg
}

func BytesToHttpRequestMessage(rawBytes []byte) (*HttpRequestMessage, error) {
	if len(rawBytes) < 4 {
		return nil, fmt.Errorf("HTTP request message is truncated")
	}
	urlLength := binary.BigEndian.Uint32(rawBytes[0:4])
	urlOffset := uint64(4) + uint64(urlLength)
	if urlOffset+4 > uint64(len(rawBytes)) {
		return nil, fmt.Errorf("HTTP request URL length is invalid")
	}
	url := rawBytes[4:urlOffset]

	dataLength := binary.BigEndian.Uint32(rawBytes[urlOffset : urlOffset+4])
	dataOffset := urlOffset + 4 + uint64(dataLength)
	if dataOffset != uint64(len(rawBytes)) {
		return nil, fmt.Errorf("HTTP request body length is invalid")
	}
	data := rawBytes[urlOffset+4 : dataOffset]

	return &HttpRequestMessage{
		UrlLength:  urlLength,
		Url:        string(url),
		DataLength: dataLength,
		Data:       data,
	}, nil
}

// Authentication messge contains the token to ensure the client
// is connected.
type AuthMessage struct {
	TokenLength uint32

	Token string
}

func NewAuthMessage(token string) AuthMessage {
	return AuthMessage{
		TokenLength: uint32(len(token)),
		Token:       token,
	}
}

func (am *AuthMessage) Bytes() []byte {
	msg := make([]byte, 0)

	msg = binary.BigEndian.AppendUint32(msg, am.TokenLength)

	msg = append(msg, []byte(am.Token)...)

	return msg
}

func bytesToAuthMessage(rawBytes []byte) (*AuthMessage, error) {
	if len(rawBytes) < 4 {
		return nil, fmt.Errorf("authentication message is truncated")
	}
	tokenLength := binary.BigEndian.Uint32(rawBytes[0:4])

	offset := uint64(4) + uint64(tokenLength)
	if offset != uint64(len(rawBytes)) {
		return nil, fmt.Errorf("authentication token length is invalid")
	}

	token := rawBytes[4:offset]

	return &AuthMessage{
		tokenLength,

		string(token),
	}, nil
}
