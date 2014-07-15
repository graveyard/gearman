package job

type stateType struct {
	i int
}

type states struct {
	Running   stateType
	Completed stateType
	Failed    stateType
}

var State = &states{stateType{0}, stateType{1}, stateType{2}}

type Job interface {
	Handle() string
	Data() chan []byte
	Warnings() chan []byte
	Status() int
	State() stateType
	SetState(stateType)
}

type job struct {
	handle         string
	data, warnings chan []byte
	status         int
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

func (j *job) Status() int {
	return j.status
}

func (j *job) State() stateType {
	return j.state
}

func (j *job) SetState(state stateType) {
	j.state = state
}

func New(string handle) Job {
	j := &job{handle: handle}
	j.data = make(chan []byte)
	j.warnings = make(chan []byte)
	j.state = State.Running
	return j
}
