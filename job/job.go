package job

type stateType struct {
	i int
}

type states struct {
	Running   stateType
	Completed stateType
	Failed    stateType
}

// Status of a Gearman job
type Status struct {
	Numerator   int32
	Denominator int32
}

// States of a Gearman job: State.Running, State.Completed, State.Failed
var State = &states{stateType{0}, stateType{1}, stateType{2}}

// Job represents a Gearman job
type Job interface {
	// The handle of the job
	Handle() string
	// Data returns a channel of work data sent by the job
	Data() chan []byte
	// Warnings returns a channel of warnings sent by the job
	Warnings() chan []byte
	// Status returns the current status of the gearman job
	Status() *Status
	// State returns the current state of the gearman job.
	// One of: State.Running, State.Completed, or State.Failed
	State() stateType
	// Sets the state for the job
	SetState(stateType)
}

type job struct {
	handle         string
	data, warnings chan []byte
	status         *Status
	state          stateType
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

func (j *job) State() stateType {
	return j.state
}

func (j *job) SetState(state stateType) {
	j.state = state
}

// New creates a new Gearman job with the specified handle
func New(handle string) Job {
	j := &job{handle: handle}
	j.data = make(chan []byte)
	j.warnings = make(chan []byte)
	j.state = State.Running
	j.status = &Status{0, 0}
	return j
}
