package gearman

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Clever/gearman/job"
	"github.com/Clever/gearman/packet"
	"github.com/Clever/gearman/scanner"
	"io"
	"net"
	"sync"
)

type Client interface {
	Close() error
	Submit(fn string, data []byte) (job.Job, error)
}

type client struct {
	conn    io.WriteCloser
	packets chan *packet.Packet
	jobs    map[string]job.Job
	newJobs chan job.Job
	jobLock sync.RWMutex
}

func (c *client) Close() error {
	c.conn.Close()
	// TODO: figure out when to close packet chan
	return nil
}

func (c *client) Submit(fn string, data []byte) (job.Job, error) {
	code := []byte{0}
	code = append(code, []byte("REQ")...)
	pack := &packet.Packet{Code: code, Type: packet.SubmitJob, Arguments: [][]byte{[]byte(fn), []byte{}, data}}
	bytes, err := pack.Bytes()
	if err != nil {
		return nil, err
	}
	written := 0
	for written != len(bytes) {
		n, err := c.conn.Write(bytes[written:len(bytes)])
		if err != nil {
			return nil, err
		}
		written += n
	}
	return <-c.newJobs, nil
}

func (c *client) addJob(j job.Job) {
	c.jobLock.Lock()
	defer c.jobLock.Unlock()
	c.jobs[j.Handle()] = j
}

func (c *client) getJob(handle string) job.Job {
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
		pack, err := packet.New(scanner.Bytes())
		if err != nil {
			fmt.Printf("ERROR PARSING PACKET! %#v\n", err)
		} else {
			c.packets <- pack
		}
	}
	if scanner.Err() != nil {
		fmt.Printf("ERROR SCANNING! %#v\n", scanner.Err())
	}
}

func (c *client) handlePackets() {
	for pack := range c.packets {
		switch pack.Type {
		case packet.JobCreated:
			j := job.New(pack.Handle())
			c.addJob(j)
			c.newJobs <- j
		case packet.WorkStatus:
			j := c.getJob(pack.Handle())
			if err := binary.Read(bytes.NewBuffer(pack.Arguments[1]), binary.BigEndian, &j.Status().Numerator); err != nil {
				fmt.Println("Error decoding numerator", err)
			}
			if err := binary.Read(bytes.NewBuffer(pack.Arguments[2]), binary.BigEndian, &j.Status().Denominator); err != nil {
				fmt.Println("Error decoding denominator", err)
			}
		case packet.WorkComplete:
			j := c.getJob(pack.Handle())
			j.SetState(job.State.Completed)
			close(j.Data())
			close(j.Warnings())
			c.deleteJob(pack.Handle())
		case packet.WorkFail:
			j := c.getJob(pack.Handle())
			j.SetState(job.State.Failed)
			close(j.Data())
			close(j.Warnings())
			c.deleteJob(pack.Handle())
		case packet.WorkData:
			j := c.getJob(pack.Handle())
			j.Data() <- pack.Arguments[1]
		case packet.WorkWarning:
			j := c.getJob(pack.Handle())
			j.Warnings() <- pack.Arguments[1]
		default:
			fmt.Println("WARNING: Unimplemented packet type", pack.Type)
		}
	}
}

func NewClient(network, addr string) (Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	c := &client{
		conn:    conn,
		packets: make(chan *packet.Packet),
		newJobs: make(chan job.Job),
		jobs:    make(map[string]job.Job),
	}
	go c.read(scanner.New(conn))

	go c.handlePackets()

	return c, nil
}
