package job

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/Clever/gearman.v2/packet"
	"gopkg.in/Clever/gearman.v2/utils"
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
		Type:      packet.Type(kind),
		Arguments: arguments,
	}
}

func TestHandlePacketsComplete(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", nil, nil, packets)
	packets <- handlePacket("", packet.WorkComplete, nil)
	assert.Equal(t, j.Run(), Completed)
}

func TestHandlePacketsFailed(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", nil, nil, packets)
	packets <- handlePacket("", packet.WorkFail, nil)
	assert.Equal(t, j.Run(), Failed)
}

func TestHandlePacketsStatus(t *testing.T) {
	packets := make(chan *packet.Packet)
	j := New("0", nil, nil, packets)
	packets <- statusPacket("", 10, 100)
	packets <- handlePacket("", packet.WorkComplete, nil)
	j.Run()
	assert.Equal(t, j.Status().Numerator, 10)
	assert.Equal(t, j.Status().Denominator, 100)
}

func TestHandlePacketsDataWarning(t *testing.T) {
	packets := make(chan *packet.Packet)
	data := utils.NewBuffer()
	warnings := utils.NewBuffer()
	j := New("0", data, warnings, packets)
	packets <- handlePacket("", packet.WorkData, [][]byte{[]byte("some data")})
	packets <- handlePacket("", packet.WorkWarning, [][]byte{[]byte("some warning")})
	packets <- handlePacket("", packet.WorkComplete, nil)
	j.Run()
	assert.Equal(t, data.Bytes(), []byte("some data"))
	assert.Equal(t, warnings.Bytes(), []byte("some warning"))
}

func TestStateStringer(t *testing.T) {
	assert.Equal(t, "Running", fmt.Sprintf("%s", Running))
	assert.Equal(t, "Completed", Completed.String())
}
