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

func (packet *gearmanPacket) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(packet.code)
	if err := binary.Write(buf, binary.BigEndian, packetType); err != nil {
		return nil, err
	}
	size := len(arguments) - 1 // One for each null-byte separator
	for _, argument := range arguments {
		size += len(argument)
	}
	if err := binary.Write(buf, binary.BigEndian, size); err != nil {
		return nil, err
	}
	// Need special handling for last argument (don't write null byte)
	for _, argument := range arguments[0 : len(arguments)-1] {
		buffer.Write(argument)
		buffer.Write([]byte{0})
	}
	buffer.Write(arguments[len(arguments)-1])
	return buffer.Bytes(), nil
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
	code := []byte{0}
	code = append(code, []byte("REQ"))
	packet := gearmanPacket{code: code, packetType: 7}
	bytes, err := packet.Bytes()
	n, err := c.conn.Write(bytes)
	// TODO: handl when n is less than len(bytes)
	if err != nil {
		return err
	}
	// TODO: wait until we get a JOB_CREATED event to get the handle, then return the job
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
