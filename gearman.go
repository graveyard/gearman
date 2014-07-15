package gearman

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Clever/gearman/job"
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
	packets chan *gearmanPacket
	jobs    map[string]job.Job
	handles chan string
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
	packet := gearmanPacket{code: code, packetType: 7, arguments: [][]byte{[]byte(fn), []byte{}, data}}
	bytes, err := packet.Bytes()
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
	handle := <-c.handles
	return job.New(handle), nil
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
		packet, err := newPacket(scanner.Bytes())
		if err != nil {
			fmt.Printf("ERROR PARSING PACKET! %#v\n", err)
		} else {
			c.packets <- packet
		}
	}
	if scanner.Err() != nil {
		fmt.Printf("ERROR SCANNING! %#v\n", scanner.Err())
	}
}

func (c *client) handlePackets() {
	for packet := range c.packets {
		switch packet.packetType {
		case JobCreated:
			c.handles <- packet.Handle()
		case WorkStatus:
			j := c.getJob(packet.Handle())
			if err := binary.Read(bytes.NewBuffer(packet.arguments[1]), binary.BigEndian, &j.Status().Numerator); err != nil {
				fmt.Println("Error decoding numerator", err)
			}
			if err := binary.Read(bytes.NewBuffer(packet.arguments[2]), binary.BigEndian, &j.Status().Denominator); err != nil {
				fmt.Println("Error decoding denominator", err)
			}
		case WorkComplete:
			j := c.getJob(packet.Handle())
			j.SetState(job.State.Completed)
			close(j.Data())
			close(j.Warnings())
		case WorkFail:
			j := c.getJob(packet.Handle())
			j.SetState(job.State.Failed)
			close(j.Data())
			close(j.Warnings())
		case WorkData:
			j := c.getJob(packet.Handle())
			j.Data() <- packet.arguments[1]
		case WorkWarning:
			j := c.getJob(packet.Handle())
			j.Warnings() <- packet.arguments[1]
		default:
			fmt.Println("WARNING: Unimplemented packet type", packet.packetType)
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
		packets: make(chan *gearmanPacket),
		handles: make(chan string),
	}
	go c.read(scanner.New(conn))

	go c.handlePackets()

	return c, nil
}
