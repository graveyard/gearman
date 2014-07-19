# gearman
--
    import "github.com/Clever/gearman"


## Usage

#### type Client

```go
type Client interface {
	// Closes the connection to the server
	Close() error
	// Submits a new job to the server with the specified function and workload
	Submit(fn string, data []byte) (job.Job, error)
}
```

Client is a Gearman client

#### func  NewClient

```go
func NewClient(network, addr string) (Client, error)
```
NewClient returns a new Gearman client pointing at the specified server
