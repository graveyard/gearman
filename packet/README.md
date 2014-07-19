# packet
--
    import "github.com/Clever/gearman/packet"


## Usage

```go
const (
	// SubmitJob = SUBMIT_JOB
	SubmitJob = 7
	// JobCreated = JOB_CREATED
	JobCreated = 8
	// WorkStatus = WORK_STATUS
	WorkStatus = 12
	// WorkComplete = WORK_COMPLETE
	WorkComplete = 13
	// WorkFail = WORK_FAIL
	WorkFail = 14
	// WorkData = WORK_DATA
	WorkData = 28
	// WorkWarning = WORK_WARNING
	WorkWarning = 29
)
```

#### type Packet

```go
type Packet struct {
	// The Code for the packet: either \0REQ or \0RES
	Code []byte
	// The Type of the packet, e.g. WorkStatus
	Type int
	// The Arguments of the packet
	Arguments [][]byte
}
```

Packet contains a Gearman packet. See http://gearman.org/protocol/

#### func  New

```go
func New(data []byte) (*Packet, error)
```
New constructs a new Packet from the slice of bytes

#### func (*Packet) Bytes

```go
func (packet *Packet) Bytes() ([]byte, error)
```
Bytes encodes the Packet into a slice of Bytes

#### func (*Packet) Handle

```go
func (packet *Packet) Handle() string
```
Handle assumes that the first argument of the packet is the job handle, returns
it as a string
