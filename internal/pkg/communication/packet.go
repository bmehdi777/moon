package communication

import (
	"encoding/binary"
	"fmt"
)

const VERSION uint8 = 1

var PacketVersionIncompatible = fmt.Errorf("Packet version are incompatible - current version : %d", VERSION)

type MessageType uint8

const (
	ConnectionStart MessageType = iota
	ConnectionClose
	HttpRequest
)

type Packet struct {
	Version  uint8
	Type     MessageType
	LenToken uint32
	LenData  uint32
	Token    string
	Data     []byte
}

func NewPacket(msgType MessageType, token string, data []byte) *Packet {
	packet := Packet{
		Version:  VERSION,
		Type:     msgType,
		LenToken: uint32(len(token)),
		LenData:  uint32(len(data)),
		Token:    token,
		Data:     data,
	}

	return &packet
}

func (p *Packet) Bytes() []byte {
	buffer := []byte{
		VERSION,
		byte(p.Type),
	}

	buffer = binary.BigEndian.AppendUint32(buffer, p.LenToken)
	buffer = binary.BigEndian.AppendUint32(buffer, p.LenData)

	buffer = append(buffer, []byte(p.Token)...)
	buffer = append(buffer, p.Data...)

	return buffer
}

func PacketFromBytes(data []byte) (*Packet, error) {
	version := uint8(data[0])
	msgType := MessageType(data[1])
	lenToken := binary.BigEndian.Uint32(data[2:6])
	lenData := binary.BigEndian.Uint32(data[6:10])

	var tokenOffset uint32 = 10
	var token string
	if lenToken != 0 {
		tokenOffset = 10 + lenToken
		token = string(data[10:tokenOffset])
	}

	var payload []byte
	if lenData != 0 {
		payload = data[tokenOffset+1:]
	}

	if version != VERSION {
		return nil, PacketVersionIncompatible
	}

	packet := Packet{
		version,
		msgType,
		lenToken,
		lenData,
		token,
		payload,
	}
	return &packet, nil
}
