# gearman
--
    import "github.com/Clever/gearman"

Package gearman provides a thread-safe Gearman client


### Example

Here's an example program that submits a job to Gearman and listens for events
from that job:

    package main

    import(
    	"github.com/Clever/gearman"
    	"github.com/Clever/gearman/job"
    )

    func main() {
    	client, err := gearman.NewClient("tcp4", "localhost:4730")
    	if err != nil {
    		panic(err)
    	}

    	j, err := client.Submit("reverse", []byte("hello world!"))
    	if err != nil {
    		panic(err)
    	}
    	// Warnings are impossible for this worker, so we don't need to listen for warnings.
    	// If they were and we didn't listen for warnings, we would block parsing from the server.
    	for data := range j.Data() {
    		println(data) // !dlrow olleh
    	}
    	println(j.State()) // job.State.Completed
    }

## Usage

#### type Client

```go
type Client interface {
	// Closes the connection to the server
	Close() error
	// Submits a new job to the server with the specified function and payload.
	// You can optionally provide custom ReadWriters for work data and warnings to be written to.
	// Otherwise, it's buffered internally in the returned Job.
	Submit(fn string, payload []byte, data, warnings io.ReadWriter) (job.Job, error)
}
```

Client is a Gearman client

#### func  NewClient

```go
func NewClient(network, addr string) (Client, error)
```
NewClient returns a new Gearman client pointing at the specified server
