package job

import (
	"github.com/Clever/gearman/packet"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func statusPacket(handle string, num, den int) *packet.Packet {
	arguments := [][]byte{}
	arguments = append(arguments, []byte(strconv.Itoa(num)), []byte(strconv.Itoa(den)))
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

func TestHandlePacketsComplete(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", packets)
	packets <- handlePacket("", packet.WorkComplete, nil)
	assert.Equal(t, j.Run(), Completed)
}

func TestHandlePacketsFailed(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", packets)
	packets <- handlePacket("", packet.WorkFail, nil)
	assert.Equal(t, j.Run(), Failed)
}

func TestHandlePacketsStatus(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", packets)
	packets <- statusPacket("", 10, 100)
	packets <- handlePacket("", packet.WorkComplete, nil)
	j.Run()
	assert.Equal(t, j.Status().Numerator, 10)
	assert.Equal(t, j.Status().Denominator, 100)
}

func TestHandlePacketsDataWarning(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", packets)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		assert.Equal(t, <-j.Data(), []byte("some data"))
	}()
	go func() {
		defer wg.Done()
		assert.Equal(t, <-j.Warnings(), []byte("some warning"))
	}()
	packets <- handlePacket("", packet.WorkData, [][]byte{[]byte("some data")})
	packets <- handlePacket("", packet.WorkWarning, [][]byte{[]byte("some warning")})
	wg.Wait()
}
