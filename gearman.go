package gearman

import (
	"bufio"
	"bytes"
	"encoding/binary"
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
	if err := binary.Write(buf, binary.BigEndian, packet.packetType); err != nil {
		return nil, err
	}
	size := len(packet.arguments) - 1 // One for each null-byte separator
	for _, argument := range packet.arguments {
		size += len(argument)
	}
	if err := binary.Write(buf, binary.BigEndian, size); err != nil {
		return nil, err
	}
	// Need special handling for last argument (don't write null byte)
	for _, argument := range packet.arguments[0 : len(packet.arguments)-1] {
		buf.Write(argument)
		buf.Write([]byte{0})
	}
	buf.Write(packet.arguments[len(packet.arguments)-1])
	return buf.Bytes(), nil
}

func newPacket(data []byte) (*gearmanPacket, error) {
	packetType := 0
	if err := binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &packetType); err != nil {
		return nil, err
	}
	arguments := bytes.Split(data[12:len(data)], []byte{0})
	return &gearmanPacket{code: data[0:4], packetType: packetType, arguments: arguments}, nil
}

type client struct {
	conn    io.WriteCloser
	packets chan *gearmanPacket
	jobs    map[string]job.Job
	handles chan string
}

func (c *client) Close() error {
	c.conn.Close()
	// TODO: figure out when to close packet chan
	return nil
}

func (c *client) Submit(fn string, data []byte) (job.Job, error) {
	code := []byte{0}
	code = append(code, []byte("REQ")...)
	packet := gearmanPacket{code: code, packetType: 7}
	bytes, err := packet.Bytes()
	if err != nil {
		return nil, err
	}
	n, err := c.conn.Write(bytes)
	if err != nil {
		return nil, err
	}
	// TODO: handl when n is less than len(bytes)
	if n != len(bytes) {
		println("Didn't write all the bytes!")
	}
	handle := <-c.handles
	return job.New(handle), nil
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
		switch packet.packetType {
		case JobCreated:
			// Parse out handle and push it onto handles
		case WorkStatus:
			// Parse out handle and update the appropriate job
		case WorkComplete:
			// Parse out handle and update the appropriate job
		case WorkFail:
			// Parse out handle and update the appropriate job
		case WorkData:
			// Parse out handle and update the appropriate job
		case WorkWarning:
			// Parse out handle and update the appropriate job
		default:
			println("WARNING: Unimplmeneted packet type", packet.packetType)
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
