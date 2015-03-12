package scanner

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/Clever/gearman.v2/packet"
)

func packetWithArgs(args [][]byte) *packet.Packet {
	return &packet.Packet{
		Code:      []byte{0, 1, 2, 3},
		Type:      1,
		Arguments: args,
	}
}

func TestScanner(t *testing.T) {
	arg := []byte("arg")
	packetNoArgs := packetWithArgs([][]byte{})
	packetOneArg := packetWithArgs([][]byte{arg})
	packetMultArgs := packetWithArgs([][]byte{arg, arg, arg})
	tmp := []byte{}
	noArgB, err := packetNoArgs.MarshalBinary()
	assert.Nil(t, err, nil)
	tmp = append(tmp, noArgB...)
	oneArgB, err := packetOneArg.MarshalBinary()
	assert.Nil(t, err, nil)
	tmp = append(tmp, oneArgB...)
	multArgB, err := packetMultArgs.MarshalBinary()
	assert.Nil(t, err, nil)
	tmp = append(tmp, multArgB...)
	buf := bytes.NewBuffer(tmp)

	scanner := New(buf)
	assert.True(t, scanner.Scan())
	assert.Equal(t, scanner.Bytes(), noArgB)
	assert.True(t, scanner.Scan())
	assert.Equal(t, scanner.Bytes(), oneArgB)
	assert.True(t, scanner.Scan())
	assert.Equal(t, scanner.Bytes(), multArgB)
}
