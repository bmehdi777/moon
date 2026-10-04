package communication

import (
	"encoding/binary"
	"fmt"
)

const VERSION uint8 = 1
const HEADER_SIZE int = 10
const READ_BUFFER_SIZE int = 1024
const MAX_PACKET_SIZE uint64 = 64 * 1024 * 1024

var PacketVersionIncompatible = fmt.Errorf("packet version is incompatible - current version: %d", VERSION)
var ErrPacketTooShort = fmt.Errorf("packet is shorter than its header")
var ErrPacketPayloadLength = fmt.Errorf("packet payload length is invalid")

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
	if len(data) < HEADER_SIZE {
		return nil, ErrPacketTooShort
	}
	header, err := HeaderFromBytes(data[0:HEADER_SIZE])
	if err != nil {
		return nil, err
	}

	if header.LenData > MAX_PACKET_SIZE {
		return nil, ErrPacketPayloadLength
	}
	expectedLength := uint64(HEADER_SIZE) + header.LenData
	if uint64(len(data)) != expectedLength {
		return nil, fmt.Errorf("%w: expected %d bytes, got %d", ErrPacketPayloadLength, expectedLength, len(data))
	}
	var payload []byte
	if header.LenData > 0 {
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
	if len(data) < HEADER_SIZE {
		return Header{}, ErrPacketTooShort
	}
	version := uint8(data[0])

	if version != VERSION {
		return Header{}, PacketVersionIncompatible
	}

	msgType := MessageType(data[1])
	lenData := binary.BigEndian.Uint64(data[2:HEADER_SIZE])
	if lenData > MAX_PACKET_SIZE {
		return Header{}, ErrPacketPayloadLength
	}

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
		var err error
		msg, err = bytesToAuthMessage(p.Payload.Data)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, fmt.Errorf("Trying to convert a raw payload to an undefined message : %v.", p.Header.Type.String())
	}

	return msg, nil
}
