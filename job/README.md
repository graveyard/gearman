# job
--
    import "github.com/Clever/gearman/job"


## Usage

#### type Job

```go
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
```

Job represents a Gearman job

#### func  New

```go
func New(handle string, packets chan *packet.Packet) Job
```
New creates a new Gearman job with the specified handle, updating the job based
on the packets in the packets channel. The only packets coming down packets
should be packets for this job.

#### type State

```go
type State int
```

State of a Gearman job

```go
const (
	// Running means that the job has not yet finished
	Running State = iota
	// Completed means that the job finished successfully
	Completed
	// Failed means that the job failed
	Failed
)
```

#### type Status

```go
type Status struct {
	Numerator   int
	Denominator int
}
```

Status of a Gearman job
