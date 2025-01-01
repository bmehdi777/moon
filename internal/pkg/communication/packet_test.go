package communication_test

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/bmehdi777/moon/internal/pkg/communication"
)

func TestPacketToByte(t *testing.T) {
	var version uint8 = 1
	var token string = "Token example"
	var data []byte = []byte("GET /")
	var lenToken uint32 = uint32(len(token))
	var lenData uint32 = uint32(len(data))

	expected := []byte{
		byte(version),
		byte(communication.ConnectionStart),
	}

	expected = binary.BigEndian.AppendUint32(expected, lenToken)
	expected = binary.BigEndian.AppendUint32(expected, lenData)

	expected = append(expected, []byte(token)...)
	expected = append(expected, data...)

	p := communication.Packet{
		Version:  1,
		Type:     communication.ConnectionStart,
		LenToken: lenToken,
		LenData:  lenData,
		Token:    token,
		Data:     data,
	}

	pBytes := p.Bytes()

	if !bytes.Equal(pBytes, expected) {
		t.Fatalf("Packets to bytes conversion isn't working. Got : %d \nExpected : %d", pBytes, expected)
	}
}

func TestByteToPacket(t *testing.T) {
	expected := communication.Packet{
		Version:  1,
		Type:     communication.ConnectionStart,
		LenToken: 0,
		LenData:  0,
	}

	rawPacket := []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	p, err := communication.PacketFromBytes(rawPacket)
	if err != nil {
		t.Fatalf("An error occured while converting raw bytes to packet : %d", err)
	}

	if !reflect.DeepEqual(p, &expected) {
		t.Fatalf("Bytes to packet conversion isn't working. Got : %#v \nExpected : %#v", p, expected)
	}
}

func TestIncompatibleVersion(t *testing.T) {
	rawPacket := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	_, err := communication.PacketFromBytes(rawPacket)
	if err == nil {
		t.Fatalf("An error should have been thrown")
	}
}
