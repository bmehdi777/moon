package communication

import (
	"encoding/binary"
	"fmt"
	"log"
)

/*
* We currently have a problem with big packet :
* We have a lenData higher than uint64 ...
* Should we do multi part package ?
*/

const VERSION uint8 = 1
const HEADER_SIZE uint32 = 14

var PacketVersionIncompatible = fmt.Errorf("Packet version are incompatible - current version : %d", VERSION)

type MessageType uint8

const (
	ConnectionStart MessageType = iota
	ConnectionClose
	HttpRequest
	HttpResponse
	InvalidToken
)

type Header struct {
	Version  uint8
	Type     MessageType
	LenToken uint32
	LenData  uint64
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
			LenData:  uint64(len(data)),
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
	buffer = binary.BigEndian.AppendUint64(buffer, p.Header.LenData)

	buffer = append(buffer, []byte(p.Payload.Token)...)
	buffer = append(buffer, p.Payload.Data...)

	return buffer
}

func PacketFromBytes(data []byte) (*Packet, error) {
	header, err := HeaderFromBytes(data[0:HEADER_SIZE])
	if err != nil {
		return nil, err
	}

	var tokenOffset uint32 = HEADER_SIZE
	var token string
	if header.LenToken != 0 {
		tokenOffset += header.LenToken
		token = string(data[HEADER_SIZE:tokenOffset])
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
	lenData := binary.BigEndian.Uint64(data[6:HEADER_SIZE])

	fmt.Println("lenData : ", lenData)

	return Header{
		Version:  version,
		Type:     msgType,
		LenToken: lenToken,
		LenData:  lenData,
	}, nil
}
