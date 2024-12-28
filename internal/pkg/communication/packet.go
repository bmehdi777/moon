package communication

import (
	"encoding/binary"
	"fmt"
)

const VERSION uint8 = 1

var PacketVersionIncompatible = fmt.Errorf("Packet version are incompatible - current version : %d", VERSION)

type MessageType uint8

const (
	ConnectionOpen MessageType = iota
	ConnectionClose
)

type Packet struct {
	Version  uint8
	Type     MessageType
	LenToken uint32
	LenData  uint32
	Token    string
	Data     []byte
}

func (p *Packet) ToByte() []byte {
	buffer := []byte{
		VERSION,
		byte(p.Type),
	}

	binary.BigEndian.AppendUint32(buffer, p.LenToken)
	binary.BigEndian.AppendUint32(buffer, p.LenData)

	buffer = append(buffer, []byte(p.Token)...)
	buffer = append(buffer, p.Data...)

	return buffer
}

func PacketFromBytes(data []byte) (*Packet, error) {
	version := uint8(data[0])
	msgType := MessageType(data[1])
	lenToken := binary.BigEndian.Uint32(data[2:5])
	lenData := binary.BigEndian.Uint32(data[6:9])
	tokenOffset := 10 + lenToken - 1
	token := string(data[10:tokenOffset])
	payload := data[tokenOffset+1:]

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
