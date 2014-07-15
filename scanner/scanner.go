package scanner

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

func needMoreData() (int, []byte, error) { return 0, nil, nil }

const HeaderSize = 12

func New(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if len(data) < HeaderSize {
			needMoreData()
		}

		var size int32
		if err := binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &size); err != nil {
			return 0, nil, err
		}

		if len(data) < HeaderSize+int(size) {
			return needMoreData()
		}

		return int(size), data[0 : HeaderSize+size], nil
	})
	return scanner
}
