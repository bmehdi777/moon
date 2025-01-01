package communication

import (
	"encoding/binary"
	"fmt"
	"log"
)

const VERSION uint8 = 1
const HEADER_SIZE uint32 = 10

var PacketVersionIncompatible = fmt.Errorf("Packet version are incompatible - current version : %d", VERSION)

type MessageType uint8

const (
	ConnectionStart MessageType = iota
	ConnectionClose
	HttpRequest
	HttpResponse
)

type Header struct {
	Version  uint8
	Type     MessageType
	LenToken uint32
	LenData  uint32
}
type Payload struct {
	Token string
	Data  []byte
}

type Packet struct {
	Header  Header
	Payload Payload
}

func NewPacket(msgType MessageType, token *string, data []byte) *Packet {
	var tok string = ""
	if token != nil {
		tok = *token
	}
	packet := Packet{
		Header: Header{
			Version:  VERSION,
			Type:     msgType,
			LenToken: uint32(len(tok)),
			LenData:  uint32(len(data)),
		},
		Payload: Payload{
			Token: tok,
			Data:  data,
		},
	}

	return &packet
}

func (p *Packet) Bytes() []byte {
	buffer := []byte{
		VERSION,
		byte(p.Header.Type),
	}

	buffer = binary.BigEndian.AppendUint32(buffer, p.Header.LenToken)
	buffer = binary.BigEndian.AppendUint32(buffer, p.Header.LenData)

	buffer = append(buffer, []byte(p.Payload.Token)...)
	buffer = append(buffer, p.Payload.Data...)

	return buffer
}

func PacketFromBytes(data []byte) (*Packet, error) {

	header, err := HeaderFromBytes(data[0:10])
	if err != nil {
		return nil, err
	}

	var tokenOffset uint32 = HEADER_SIZE
	var token string
	if header.LenToken != 0 {
		tokenOffset += header.LenToken
		token = string(data[10:tokenOffset])
	}

	var payload []byte
	if header.LenData != 0 {
		payload = data[tokenOffset:]
	}

	packet := Packet{
		Header: header,
		Payload: Payload{
			token,
			payload,
		},
	}
	return &packet, nil
}

func HeaderFromBytes(data []byte) (Header, error) {
	version := uint8(data[0])

	if version != VERSION {
		log.Println("Version received : ", version)
		return Header{}, PacketVersionIncompatible
	}

	msgType := MessageType(data[1])
	lenToken := binary.BigEndian.Uint32(data[2:6])
	lenData := binary.BigEndian.Uint32(data[6:10])

	return Header{
		Version:  version,
		Type:     msgType,
		LenToken: lenToken,
		LenData:  lenData,
	}, nil
}
