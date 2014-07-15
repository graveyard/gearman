package gearman

import (
	"bytes"
	"encoding/binary"
)

type gearmanPacket struct {
	code       []byte
	packetType int
	arguments  [][]byte
}

func (packet *gearmanPacket) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(packet.code)
	if err := binary.Write(buf, binary.BigEndian, int32(packet.packetType)); err != nil {
		return nil, err
	}
	size := len(packet.arguments) - 1 // One for each null-byte separator
	for _, argument := range packet.arguments {
		size += len(argument)
	}
	if err := binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return nil, err
	}
	if len(packet.arguments) > 0 {
		// Need special handling for last argument (don't write null byte)
		for _, argument := range packet.arguments[0 : len(packet.arguments)-1] {
			buf.Write(argument)
			buf.Write([]byte{0})
		}
		buf.Write(packet.arguments[len(packet.arguments)-1])
	}
	return buf.Bytes(), nil
}

// Handle assumes that the first argument of the packet is the job handle
func (packet *gearmanPacket) Handle() string {
	return string(packet.arguments[0])
}

func newPacket(data []byte) (*gearmanPacket, error) {
	packetType := int32(0)
	if err := binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &packetType); err != nil {
		return nil, err
	}
	arguments := bytes.Split(data[12:len(data)], []byte{0})
	return &gearmanPacket{code: data[0:4], packetType: int(packetType), arguments: arguments}, nil
}
