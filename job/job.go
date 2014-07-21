package job

import (
	"fmt"
	"github.com/Clever/gearman/packet"
	"strconv"
)

// State of a Gearman job
type State int

const (
	// Running means that the job has not yet finished
	Running State = iota
	// Completed means that the job finished successfully
	Completed
	// Failed means that the job failed
	Failed
)

// Status of a Gearman job
type Status struct {
	Numerator   int
	Denominator int
}

// Job represents a Gearman job
type Job interface {
	// The handle of the job
	Handle() string
	// Data returns a channel of work data sent by the job
	// NOTE: If you don't listen to this channel, you will block parsing of new Gearman packets
	Data() <-chan []byte
	// Warnings returns a channel of warnings sent by the job
	// NOTE: If you don't listen to this channel, you will block parsing of new Gearman packets
	Warnings() <-chan []byte
	// Status returns the current status of the gearman job
	Status() Status
	// Blocks until the job completes. Returns the state, Completed or Failed.
	Run() State
}

type job struct {
	handle         string
	data, warnings chan []byte
	status         Status
	state          State
	done           chan struct{}
}

func (j *job) Handle() string {
	return j.handle
}

func (j *job) Data() <-chan []byte {
	return j.data
}

func (j *job) Warnings() <-chan []byte {
	return j.warnings
}

func (j *job) Status() Status {
	return j.status
}

func (j *job) Run() State {
	<-j.done
	return j.state
}

func (j *job) handlePackets(packets chan *packet.Packet) {
	for pack := range packets {
		switch pack.Type {
		case packet.WorkStatus:
			num, err := strconv.Atoi(string(pack.Arguments[1]))
			if err != nil {
				fmt.Println("Error converting numerator", err)
			}
			den, err := strconv.Atoi(string(pack.Arguments[2]))
			if err != nil {
				fmt.Println("Error converting denominator", err)
			}
			j.status = Status{Numerator: num, Denominator: den}
		case packet.WorkComplete:
			j.state = Completed
			close(j.data)
			close(j.warnings)
			close(j.done)
			close(packets)
		case packet.WorkFail:
			j.state = Failed
			close(j.data)
			close(j.warnings)
			close(j.done)
			close(packets)
		case packet.WorkData:
			j.data <- pack.Arguments[1]
		case packet.WorkWarning:
			j.warnings <- pack.Arguments[1]
		default:
			fmt.Println("WARNING: Unimplemented packet type", pack.Type)
		}
	}
}

// New creates a new Gearman job with the specified handle, updating the job based on the packets
// in the packets channel. The only packets coming down packets should be packets for this job.
func New(handle string, packets chan *packet.Packet) Job {
	j := &job{
		handle:   handle,
		data:     make(chan []byte),
		warnings: make(chan []byte),
		status:   Status{0, 0},
		state:    Running,
		done:     make(chan struct{}),
	}
	go j.handlePackets(packets)
	return j
}
