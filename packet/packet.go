package packet

import (
	"bytes"
	"encoding/binary"
)

type Packet struct {
	Code      []byte
	Type      int
	Arguments [][]byte
}

func (packet *Packet) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(packet.Code)
	if err := binary.Write(buf, binary.BigEndian, int32(packet.Type)); err != nil {
		return nil, err
	}
	size := len(packet.Arguments) - 1 // One for each null-byte separator
	for _, argument := range packet.Arguments {
		size += len(argument)
	}
	if err := binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return nil, err
	}
	if len(packet.Arguments) > 0 {
		// Need special handling for last argument (don't write null byte)
		for _, argument := range packet.Arguments[0 : len(packet.Arguments)-1] {
			buf.Write(argument)
			buf.Write([]byte{0})
		}
		buf.Write(packet.Arguments[len(packet.Arguments)-1])
	}
	return buf.Bytes(), nil
}

// Handle assumes that the first argument of the packet is the job handle
func (packet *Packet) Handle() string {
	return string(packet.Arguments[0])
}

func New(data []byte) (*Packet, error) {
	packetType := int32(0)
	if err := binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &packetType); err != nil {
		return nil, err
	}
	arguments := bytes.Split(data[12:len(data)], []byte{0})
	return &Packet{Code: data[0:4], Type: int(packetType), Arguments: arguments}, nil
}
