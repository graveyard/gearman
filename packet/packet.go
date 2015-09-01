// Package packet provides structures to marshal binary data to and from binary data.
// The specification is located at http://gearman.org/protocol/.
package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type packetCode []byte

var (
	// Req is the code for a Request packet
	Req = packetCode([]byte{0, byte('R'), byte('E'), byte('Q')})
	// Res is the code for a Response packet
	Res = packetCode([]byte{0, byte('R'), byte('E'), byte('S')})
)

// Packet contains a Gearman packet. See http://gearman.org/protocol/
type Packet struct {
	// The Code for the packet: either \0REQ or \0RES
	Code packetCode
	// The Type of the packet, e.g. WorkStatus
	Type Type
	// The Arguments of the packet
	Arguments [][]byte
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (packet *Packet) UnmarshalBinary(data []byte) error {
	// ensure packet is of minimum length
	if len(data) < 12 {
		return errors.New("All gearman packets must be at least 12 bytes.")
	}

	// determine the packet magic code
	if bytes.Compare(data[0:4], Req) == 0 {
		packet.Code = Req
	} else if bytes.Compare(data[0:4], Res) == 0 {
		packet.Code = Res
	} else {
		return fmt.Errorf("unrecognized packet code %#v", data[0:4])
	}

	// determine the kind of packet
	kind := int32(0)
	if err := binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &kind); err != nil {
		return err
	}
	packet.Type = Type(kind)

	// parse the length of the packet
	length := int32(0)
	if err := binary.Read(bytes.NewBuffer(data[8:12]), binary.BigEndian, &length); err != nil {
		return err
	}

	// parse the arguments into a byte array
	packet.Arguments = [][]byte{}
	if length > 0 {
		packet.Arguments = bytes.Split(data[12:len(data)], []byte{0})
	}

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (packet *Packet) MarshalBinary() ([]byte, error) {
	// form a buffer with the packet's magic code
	buf := bytes.NewBuffer(packet.Code)

	// write the request header
	if err := binary.Write(buf, binary.BigEndian, int32(packet.Type)); err != nil {
		return nil, err
	}

	// write the size of the packet
	size := len(packet.Arguments) - 1 // One for each null-byte separator
	for _, argument := range packet.Arguments {
		size += len(argument)
	}
	size = int(math.Max(0, float64(size)))

	// write the size of the packet
	if err := binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return nil, err
	}

	// write all arguments provided
	for i := 0; i < len(packet.Arguments); i++ {
		if _, err := buf.Write(packet.Arguments[i]); err != nil {
			return nil, err
		}

		// null deliminate every argument but the last
		if i == len(packet.Arguments)-1 {
			break // last argument
		} else if err := buf.WriteByte(0); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
