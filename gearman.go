package gearman

import (
	"bufio"
	"encoding/binary"
	"github.com/Clever/gearman/job"
	"github.com/Clever/gearman/scanner"
	"io"
	"net"
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

func (packet *gearmanPacket) Bytes() []byte {
	buf := bytes.NewBuffer(packet.code)
	binary.Write(buf, binary.BigEndian, packetType)
	// TODO: write size, convert arguments
	return nil
}

func newPacket(data []byte) (*gearmanPacket, error) {
	// TODO: parse bytes into packet
	return nil, nil
}

type client struct {
	conn    io.WriteCloser
	packets chan *gearmanPacket
	jobs    map[string]job.Job
}

func (c *client) Close() error {
	// TODO: close connection, figure out when to close packet chan
	return nil
}

func (c *client) Submit(fn string, data []byte) (job.Job, error) {
	// TODO
	// create a gearmanPacket, send it
	// wait until we get a JOB_CREATED event to get the handle, then return
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
		// Basically a switch on packet type, and then do something based on the arguments
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
