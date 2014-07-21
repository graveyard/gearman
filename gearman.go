package gearman

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Clever/gearman/job"
	"github.com/Clever/gearman/packet"
	"github.com/Clever/gearman/scanner"
	"io"
	"net"
	"sync"
)

// Client is a Gearman client
type Client interface {
	// Closes the connection to the server
	Close() error
	// Submits a new job to the server with the specified function and workload
	Submit(fn string, data []byte) (job.Job, error)
}

type client struct {
	conn    io.WriteCloser
	packets chan *packet.Packet
	jobs    map[string]chan *packet.Packet
	newJobs chan job.Job
	jobLock sync.RWMutex
}

func (c *client) Close() error {
	c.conn.Close()
	// TODO: figure out when to close packet chan
	return nil
}

func (c *client) Submit(fn string, data []byte) (job.Job, error) {
	pack := &packet.Packet{Code: packet.Req, Type: packet.SubmitJob, Arguments: [][]byte{[]byte(fn), []byte{}, data}}
	b, err := pack.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(c.conn, bytes.NewBuffer(b)); err != nil {
		return nil, err
	}
	return <-c.newJobs, nil
}

func (c *client) addJob(handle string, packets chan *packet.Packet) {
	c.jobLock.Lock()
	defer c.jobLock.Unlock()
	c.jobs[handle] = packets
}

func (c *client) getJob(handle string) chan *packet.Packet {
	c.jobLock.RLock()
	defer c.jobLock.RUnlock()
	return c.jobs[handle]
}

func (c *client) deleteJob(handle string) {
	c.jobLock.Lock()
	defer c.jobLock.Unlock()
	delete(c.jobs, handle)
}

func (c *client) read(scanner *bufio.Scanner) {
	for scanner.Scan() {
		pack := &packet.Packet{}
		if err := pack.UnmarshalBinary(scanner.Bytes()); err != nil {
			fmt.Printf("ERROR PARSING PACKET! %#v\n", err)
		} else {
			c.packets <- pack
		}
	}
	if scanner.Err() != nil {
		fmt.Printf("ERROR SCANNING! %#v\n", scanner.Err())
	}
}

func (c *client) routePackets() {
	for pack := range c.packets {
		handle := string(pack.Arguments[0])
		if pack.Type == packet.JobCreated {
			packets := make(chan *packet.Packet)
			j := job.New(handle, packets)
			c.addJob(handle, packets)
			c.newJobs <- j
			go func() {
				_ = j.Run()
				c.deleteJob(handle)
			}()
		} else {
			c.getJob(handle) <- pack
		}
	}
}

// NewClient returns a new Gearman client pointing at the specified server
func NewClient(network, addr string) (Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	c := &client{
		conn:    conn,
		packets: make(chan *packet.Packet),
		newJobs: make(chan job.Job),
		jobs:    make(map[string]chan *packet.Packet),
	}
	go c.read(scanner.New(conn))

	go c.routePackets()

	return c, nil
}
