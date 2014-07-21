package job

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
	Numerator   int32
	Denominator int32
}

// Job represents a Gearman job
type Job interface {
	// The handle of the job
	Handle() string
	// Data returns a channel of work data sent by the job
	// NOTE: If you don't listen to this channel, you will block parsing of new Gearman packets
	Data() chan []byte
	// Warnings returns a channel of warnings sent by the job
	// NOTE: If you don't listen to this channel, you will block parsing of new Gearman packets
	Warnings() chan []byte
	// Status returns the current status of the gearman job
	Status() *Status
	// State returns the current state of the gearman job.
	// One of: State.Running, State.Completed, or State.Failed
	State() State
	// Sets the state for the job
	SetState(State)
}

type job struct {
	handle         string
	data, warnings chan []byte
	status         *Status
	state          State
}

func (j *job) Handle() string {
	return j.handle
}

func (j *job) Data() chan []byte {
	return j.data
}

func (j *job) Warnings() chan []byte {
	return j.warnings
}

func (j *job) Status() *Status {
	return j.status
}

func (j *job) State() State {
	return j.state
}

func (j *job) SetState(state State) {
	j.state = state
}

// New creates a new Gearman job with the specified handle
func New(handle string) Job {
	return &job{
		handle:   handle,
		data:     make(chan []byte),
		warnings: make(chan []byte),
		status:   &Status{0, 0},
		state:    Running,
	}
}
