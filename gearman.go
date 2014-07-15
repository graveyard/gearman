package gearman

import (
	"bufio"
	"github.com/Clever/gearman/job"
	"github.com/Clever/gearman/scanner"
	"io"
	"net"
)

type Client interface {
	Close() error
	Submit(fn string, data []byte) (job.Job, error)
}

type gearmanPacket struct{}

func newPacket(data []byte) (*gearmanPacket, error) {
	// TODO
	return nil, nil
}

type client struct {
	conn    io.WriteCloser
	packets chan *gearmanPacket
	jobs    map[string]job.Job
}

func (c *client) Close() error {
	// TODO
	return nil
}

func (c *client) Submit(fn string, data []byte) (job.Job, error) {
	// TODO
	return nil, nil
}

func (c *client) read(scanner *bufio.Scanner) {
	for scanner.Scan() {
		packet, err := newPacket(scanner.Bytes())
		if err != nil {
			println("ERROR PARSING PACKET!")
		}
		c.packets <- packet
	}
	if scanner.Err() != nil {
		println("ERROR SCANNING!")
	}
}

func (c *client) handlePackets() {
	for packet := range c.packets {
		// Basically a giant switch on packet type, and then do something based on the handle
	}
}

func NewClient(network, addr string) (Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return err
	}
	c := &client{
		conn:    conn,
		packets: make(chan *gearmanPacket),
	}
	go read(scanner.New(conn))

	for i := 0; i < 100; i++ {
		go handlePackets()
	}

	return c, nil
}
