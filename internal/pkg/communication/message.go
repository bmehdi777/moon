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
	case InvalidToken:
		return "InvalidToken"
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

	InvalidToken
)

// Complexe payload are converted to type
type Message interface {
	Bytes() []byte
}

// Authentication messge contains the token to ensure the client
// is connected.
type AuthMessage struct {
	TokenLength       uint32

	Token       string
}

func NewAuthMessage(token string) AuthMessage {
	return AuthMessage{
		TokenLength:       uint32(len(token)),
		Token:             token,
	}
}

func (am *AuthMessage) Bytes() []byte {
	msg := make([]byte, 0)

	msg = binary.BigEndian.AppendUint32(msg, am.TokenLength)

	msg = append(msg, []byte(am.Token)...)

	return msg
}
func bytesToAuthMessage(rawBytes []byte) *AuthMessage {
	tokenLength := binary.BigEndian.Uint32(rawBytes[0:4])

	offset := 4 + tokenLength

	token := rawBytes[4:offset]

	return &AuthMessage{
		tokenLength,

		string(token),
	}
}
