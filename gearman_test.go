package gearman

import (
	"bytes"
	"encoding/binary"
	"github.com/Clever/gearman/job"
	"github.com/Clever/gearman/packet"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

type bufferCloser struct {
	bytes.Buffer
}

func (buf *bufferCloser) Close() error {
	return nil
}

func mockClient() *client {
	c := &client{
		conn:    &bufferCloser{},
		packets: make(chan *packet.Packet),
		// Add buffers to prevent blocking in test cases
		newJobs: make(chan job.Job, 10),
		jobs:    make(map[string]job.Job, 10),
	}
	go c.handlePackets()
	return c
}

func TestSubmit(t *testing.T) {
	c := mockClient()
	buf := c.conn.(*bufferCloser)
	expected := job.New("the_handle")
	c.newJobs <- expected
	j, err := c.Submit("my_function", []byte("my data"))
	assert.Nil(t, err)
	assert.Equal(t, j, expected)
	expectedPacket := &packet.Packet{
		Code:      []byte{0x0, 0x52, 0x45, 0x51}, // \0REQ
		Type:      packet.SubmitJob,
		Arguments: [][]byte{[]byte("my_function"), []byte{}, []byte("my data")},
	}
	b, err := expectedPacket.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, buf.Bytes(), b)
}

func statusPacket(handle string, numerator, denominator int32) *packet.Packet {
	arguments := [][]byte{}
	bufNum := bytes.NewBuffer(make([]byte, 4))
	bufDen := bytes.NewBuffer(make([]byte, 4))
	if err := binary.Write(bufNum, binary.BigEndian, numerator); err != nil {
		panic(err)
	}
	if err := binary.Write(bufDen, binary.BigEndian, denominator); err != nil {
		panic(err)
	}
	arguments = append(arguments, bufNum.Bytes(), bufDen.Bytes())
	return handlePacket(handle, packet.WorkStatus, arguments)
}

func handlePacket(handle string, kind int, arguments [][]byte) *packet.Packet {
	if arguments == nil {
		arguments = [][]byte{}
	}
	arguments = append([][]byte{[]byte(handle)}, arguments...)
	return &packet.Packet{
		Type:      packet.PacketType(kind),
		Arguments: arguments,
	}
}

func TestHandlePackets(t *testing.T) {
	c := mockClient()
	for i := 0; i < 5; i++ {
		c.jobs[strconv.Itoa(i)] = job.New(strconv.Itoa(i))
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		for data := range c.jobs["3"].Data() {
			assert.Equal(t, data, []byte("data!!"))
		}
	}()
	go func() {
		defer wg.Done()
		for warning := range c.jobs["4"].Warnings() {
			assert.Equal(t, warning, []byte("warning :(:("))
		}
	}()
	go func() {
		defer wg.Done()
		j := <-c.newJobs
		assert.Equal(t, j.Handle(), "5")
	}()
	c.packets <- statusPacket("0", 10, 100)
	c.packets <- handlePacket("1", packet.WorkComplete, nil)
	c.packets <- handlePacket("2", packet.WorkFail, nil)
	c.packets <- handlePacket("3", packet.WorkData, [][]byte{[]byte("data!!")})
	c.packets <- handlePacket("4", packet.WorkWarning, [][]byte{[]byte("warning :(:(")})
	assert.Nil(t, c.jobs["1"])
	assert.Nil(t, c.jobs["2"])
	c.packets <- handlePacket("3", packet.WorkComplete, nil)
	c.packets <- handlePacket("4", packet.WorkComplete, nil)
	c.packets <- handlePacket("5", packet.JobCreated, nil)
	wg.Wait()
}
