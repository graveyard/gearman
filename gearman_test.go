package gearman

import (
	"bytes"
	"github.com/Clever/gearman/job"
	"github.com/Clever/gearman/packet"
	"github.com/stretchr/testify/assert"
	"testing"
)

type bufferCloser struct {
	bytes.Buffer
}

func (buf *bufferCloser) Close() error {
	return nil
}

func mockClient() *client {
	return &client{
		conn:    &bufferCloser{},
		packets: make(chan *packet.Packet),
		newJobs: make(chan job.Job),
		jobs:    make(map[string]job.Job),
	}
}

func TestSubmit(t *testing.T) {
	c := mockClient()
	buf := c.conn.(*bufferCloser)
	expected := job.New("the_handle")
	go func() {
		c.newJobs <- expected
	}()
	j, err := c.Submit("my_function", []byte("my data"))
	assert.Nil(t, err)
	assert.Equal(t, j, expected)
	expectedPacket := &packet.Packet{
		Code:      []byte{0x0, 0x52, 0x45, 0x51}, // \0REQ
		Type:      packet.SubmitJob,
		Arguments: [][]byte{[]byte("my_function"), []byte{}, []byte("my data")},
	}
	b, err := expectedPacket.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf.Bytes(), b)
}
