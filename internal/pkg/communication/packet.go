package communication

import (
	"encoding/binary"
	"fmt"
	"log"
)

const VERSION uint8 = 1
const HEADER_SIZE int = 10
const READ_BUFFER_SIZE int = 1024

var PacketVersionIncompatible = fmt.Errorf("Packet version are incompatible - current version : %d", VERSION)


type Header struct {
	Version uint8       // 1
	Type    MessageType // 1
	LenData uint64      // 8
}
type Payload struct {
	Data []byte
}

type Packet struct {
	Header  Header
	Payload Payload
}

func NewPacket(msgType MessageType, data []byte) *Packet {
	packet := Packet{
		Header: Header{
			Version: VERSION,
			Type:    msgType,
			LenData: uint64(len(data)),
		},
		Payload: Payload{
			Data: data,
		},
	}

	return &packet
}

func (p *Packet) Bytes() []byte {
	buffer := []byte{
		VERSION,
		byte(p.Header.Type),
	}

	buffer = binary.BigEndian.AppendUint64(buffer, p.Header.LenData)

	buffer = append(buffer, p.Payload.Data...)

	return buffer
}

func PacketFromBytes(data []byte) (*Packet, error) {
	header, err := HeaderFromBytes(data[0:HEADER_SIZE])
	if err != nil {
		return nil, err
	}

	var payload []byte
	if header.LenData != 0 {
		payload = data[HEADER_SIZE:]
	}

	packet := Packet{
		Header: header,
		Payload: Payload{
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
	lenData := binary.BigEndian.Uint64(data[2:HEADER_SIZE])

	return Header{
		Version: version,
		Type:    msgType,
		LenData: lenData,
	}, nil
}

func (p *Packet) Message() (Message, error) {
	var msg Message

	switch p.Header.Type {
	case ConnectionStart:
		msg = bytesToAuthMessage(p.Payload.Data)
		break
	default:
		return nil, fmt.Errorf("Trying to convert a raw payload to an undefined message : %v.", p.Header.Type.String())
	}

	return msg, nil
}
