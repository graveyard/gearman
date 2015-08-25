package job

import (
	"fmt"
	"io"
	"strconv"

	"gopkg.in/Clever/gearman.v2/packet"
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
type Job struct {
	handle         string
	data, warnings io.WriteCloser
	status         Status
	state          State
	done           chan struct{}
}

// Handle returns job handle
func (j Job) Handle() string {
	return j.handle
}

// Status returns the current status of the gearman job
func (j Job) Status() Status {
	return j.status
}

// Run blocks until the job completes. Returns the state, Completed or Failed.
func (j *Job) Run() State {
	<-j.done
	return j.state
}

func (j *Job) handlePackets(packets chan *packet.Packet) {
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
			close(j.done)
		case packet.WorkFail:
			j.state = Failed
			close(j.done)
		case packet.WorkData:
			if _, err := j.data.Write(pack.Arguments[1]); err != nil {
				fmt.Printf("Error writing data, arg: %s, err: %s", pack.Arguments[1], err)
			}
		case packet.WorkWarning:
			if _, err := j.warnings.Write(pack.Arguments[1]); err != nil {
				fmt.Printf("Error writing warnings, arg: %s, err: %s", pack.Arguments[1], err)
			}
		default:
			fmt.Println("WARNING: Unimplemented packet type", pack.Type)
		}
	}
}

// New creates a new Gearman job with the specified handle, updating the job based on the packets
// in the packets channel. The only packets coming down packets should be packets for this job.
// It also takes in two WriteClosers to right job data and warnings to.
func New(handle string, data, warnings io.WriteCloser, packets chan *packet.Packet) *Job {
	j := &Job{
		handle:   handle,
		data:     data,
		warnings: warnings,
		status:   Status{0, 0},
		state:    Running,
		done:     make(chan struct{}),
	}
	go j.handlePackets(packets)
	return j
}
