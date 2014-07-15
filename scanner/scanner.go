package scanner

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

func needMoreData() (int, []byte, error) { return 0, nil, nil }

const headerSize = 12

func New(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if len(data) < headerSize {
			return needMoreData()
		}

		var size int32
		if err := binary.Read(bytes.NewBuffer(data[8:12]), binary.BigEndian, &size); err != nil {
			return 0, nil, err
		}

		if len(data) < headerSize+int(size) {
			return needMoreData()
		}

		return int(headerSize + size), data[0 : headerSize+size], nil
	})
	return scanner
}
