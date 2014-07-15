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

type gearmanPacket struct {
	code       []byte
	packetType int
	arguments  [][]byte
}

func (packet *gearmanPacket) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(packet.code)
	if err := binary.Write(buf, binary.BigEndian, int32(packet.packetType)); err != nil {
		return nil, err
	}
	size := len(packet.arguments) - 1 // One for each null-byte separator
	for _, argument := range packet.arguments {
		size += len(argument)
	}
	if err := binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return nil, err
	}
	if len(packet.arguments) > 0 {
		// Need special handling for last argument (don't write null byte)
		for _, argument := range packet.arguments[0 : len(packet.arguments)-1] {
			buf.Write(argument)
			buf.Write([]byte{0})
		}
		buf.Write(packet.arguments[len(packet.arguments)-1])
	}
	return buf.Bytes(), nil
}

// Handle assumes that the first argument of the packet is the job handle
func (packet *gearmanPacket) Handle() string {
	return string(packet.arguments[0])
}

func newPacket(data []byte) (*gearmanPacket, error) {
	packetType := int32(0)
	if err := binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &packetType); err != nil {
		return nil, err
	}
	arguments := bytes.Split(data[12:len(data)], []byte{0})
	return &gearmanPacket{code: data[0:4], packetType: int(packetType), arguments: arguments}, nil
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
			_ = j
		case WorkComplete:
			j := c.getJob(packet.Handle())
			j.SetState(job.State.Completed)
		case WorkFail:
			j := c.getJob(packet.Handle())
			j.SetState(job.State.Failed)
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

	for i := 0; i < 100; i++ {
		go c.handlePackets()
	}

	return c, nil
}
