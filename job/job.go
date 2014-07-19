package job

type stateType struct {
	i int
}

type states struct {
	Running   stateType
	Completed stateType
	Failed    stateType
}

type Status struct {
	Numerator   int32
	Denominator int32
}

var State = &states{stateType{0}, stateType{1}, stateType{2}}

type Job interface {
	Handle() string
	Data() chan []byte
	Warnings() chan []byte
	Status() *Status
	State() stateType
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

func New(handle string) Job {
	j := &job{handle: handle}
	j.data = make(chan []byte)
	j.warnings = make(chan []byte)
	j.state = State.Running
	j.status = &Status{0, 0}
	return j
}
