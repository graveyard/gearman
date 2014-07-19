# job
--
    import "github.com/Clever/gearman/job"


## Usage

```go
var State = &states{stateType{0}, stateType{1}, stateType{2}}
```
States of a Gearman job: State.Running, State.Completed, State.Failed

#### type Job

```go
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
	State() stateType
	// Sets the state for the job
	SetState(stateType)
}
```

Job represents a Gearman job

#### func  New

```go
func New(handle string) Job
```
New creates a new Gearman job with the specified handle

#### type Status

```go
type Status struct {
	Numerator   int32
	Denominator int32
}
```

Status of a Gearman job
