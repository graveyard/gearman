package gearman

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/Clever/gearman.v1/job"
	"gopkg.in/Clever/gearman.v1/packet"
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
		newJobs:     make(chan job.Job, 10),
		jobs:        make(map[string]chan *packet.Packet, 10),
		partialJobs: make(chan *partialJob, 10),
	}
	go c.routePackets()
	return c
}

func TestSubmit(t *testing.T) {
	c := mockClient()
	buf := c.conn.(*bufferCloser)
	expected := job.New("the_handle", nil, nil, make(chan *packet.Packet))
	c.newJobs <- expected
	j, err := c.Submit("my_function", []byte("my data"), nil, nil)
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

func handlePacket(handle string, kind int, arguments [][]byte) *packet.Packet {
	if arguments == nil {
		arguments = [][]byte{}
	}
	arguments = append([][]byte{[]byte(handle)}, arguments...)
	return &packet.Packet{
		Type:      packet.Type(kind),
		Arguments: arguments,
	}
}

func TestJobCreated(t *testing.T) {
	c := mockClient()
	c.partialJobs <- &partialJob{nil, nil}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var j job.Job
	var packets chan *packet.Packet
	go func() {
		defer wg.Done()
		j = <-c.newJobs
		packets = c.getJob("5")
		assert.Equal(t, j.Handle(), "5")
	}()
	c.packets <- handlePacket("5", packet.JobCreated, nil)
	wg.Wait()
	c.packets <- handlePacket("5", packet.WorkComplete, nil)
	j.Run()
	<-packets // Wait until packet channel is closed, so we know that we've deleted the job
	assert.Nil(t, c.jobs["5"])
}

func TestRoutePackets(t *testing.T) {
	c := mockClient()
	packetChans := []chan *packet.Packet{}
	for i := 0; i < 5; i++ {
		packetChans = append(packetChans, make(chan *packet.Packet, 10))
		c.jobs[strconv.Itoa(i)] = packetChans[i]
	}

	packets := []*packet.Packet{}
	packets = append(packets, handlePacket("0", packet.WorkFail, nil))
	packets = append(packets, handlePacket("1", packet.WorkFail, nil))
	packets = append(packets, handlePacket("2", packet.WorkFail, nil))
	packets = append(packets, handlePacket("3", packet.WorkFail, nil))
	packets = append(packets, handlePacket("4", packet.WorkFail, nil))
	for _, pack := range packets {
		c.packets <- pack
	}
	for i := 0; i < 5; i++ {
		pack := <-packetChans[i]
		assert.Equal(t, pack, packets[i])
	}
}
