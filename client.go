package gearman

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"

	"gopkg.in/Clever/gearman.v2/job"
	"gopkg.in/Clever/gearman.v2/packet"
	"gopkg.in/Clever/gearman.v2/scanner"
)

// noOpCloser is like an ioutil.NopCloser, but for an io.Writer.
type noOpCloser struct {
	w io.Writer
}

func (c noOpCloser) Write(data []byte) (n int, err error) {
	return c.w.Write(data)
}

func (c noOpCloser) Close() error {
	return nil
}

var discard = noOpCloser{w: ioutil.Discard}

// Client is a Gearman client
type Client struct {
	conn        io.WriteCloser
	packets     chan *packet.Packet
	jobs        map[string]chan *packet.Packet
	partialJobs chan *partialJob
	newJobs     chan *job.Job
	jobLock     sync.RWMutex
}

type partialJob struct {
	data, warnings io.WriteCloser
}

// Close terminates the connection to the server
func (c *Client) Close() error {
	c.conn.Close()
	// TODO: figure out when to close packet chan
	return nil
}

func (c *Client) submit(fn string, payload []byte, data, warnings io.WriteCloser, t packet.Type) (*job.Job, error) {
	pack := &packet.Packet{
		Code:      packet.Req,
		Type:      t,
		Arguments: [][]byte{[]byte(fn), []byte{}, payload},
	}
	b, err := pack.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(c.conn, bytes.NewBuffer(b)); err != nil {
		return nil, err
	}
	c.partialJobs <- &partialJob{data: data, warnings: warnings}
	return <-c.newJobs, nil
}

// Submit sends a new job to the server with the specified function and payload. You must provide
// two WriteClosers for data and warnings to be written to.
func (c *Client) Submit(fn string, payload []byte, data, warnings io.WriteCloser) (*job.Job, error) {
	return c.submit(fn, payload, data, warnings, packet.SubmitJob)
}

// SubmitBackground submits a background job. There is no access to data, warnings, or completion
// state.
func (c *Client) SubmitBackground(fn string, payload []byte) error {
	_, err := c.submit(fn, payload, discard, discard, packet.SubmitJobBg)
	return err
}

func (c *Client) addJob(handle string, packets chan *packet.Packet) {
	c.jobLock.Lock()
	defer c.jobLock.Unlock()
	c.jobs[handle] = packets
}

func (c *Client) getJob(handle string) chan *packet.Packet {
	c.jobLock.RLock()
	defer c.jobLock.RUnlock()
	return c.jobs[handle]
}

func (c *Client) deleteJob(handle string) {
	c.jobLock.Lock()
	defer c.jobLock.Unlock()
	delete(c.jobs, handle)
}

func (c *Client) read(scanner *bufio.Scanner) {
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

func (c *Client) routePackets() {
	for pack := range c.packets {
		handle := string(pack.Arguments[0])
		if pack.Type == packet.JobCreated {
			packets := make(chan *packet.Packet)
			pj := <-c.partialJobs
			j := job.New(handle, pj.data, pj.warnings, packets)
			c.addJob(handle, packets)
			c.newJobs <- j
			go func() {
				defer close(packets)
				defer c.deleteJob(handle)
				j.Run()
			}()
		} else {
			c.getJob(handle) <- pack
		}
	}
}

// NewClient returns a new Gearman client pointing at the specified server
func NewClient(network, addr string) (*Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn:        conn,
		packets:     make(chan *packet.Packet),
		newJobs:     make(chan *job.Job),
		partialJobs: make(chan *partialJob),
		jobs:        make(map[string]chan *packet.Packet),
	}
	go c.read(scanner.New(conn))

	go c.routePackets()

	return c, nil
}
