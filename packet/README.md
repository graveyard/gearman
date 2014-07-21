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

```go
var (
	// Req is the code for a Request packet
	Req = packetCode([]byte{0, byte('R'), byte('E'), byte('Q')})
	// Res is the code for a Response packet
	Res = packetCode([]byte{0, byte('R'), byte('E'), byte('Q')})
)
```

#### type Packet

```go
type Packet struct {
	// The Code for the packet: either \0REQ or \0RES
	Code packetCode
	// The Type of the packet, e.g. WorkStatus
	Type int
	// The Arguments of the packet
	Arguments [][]byte
}
```

Packet contains a Gearman packet. See http://gearman.org/protocol/

#### func (*Packet) Handle

```go
func (packet *Packet) Handle() string
```
Handle assumes that the first argument of the packet is the job handle, returns
it as a string

#### func (*Packet) MarshalBinary

```go
func (packet *Packet) MarshalBinary() ([]byte, error)
```
MarshalBinary implements the encoding.BinaryMarshaler interface

#### func (*Packet) UnmarshalBinary

```go
func (packet *Packet) UnmarshalBinary(data []byte) error
```
UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
