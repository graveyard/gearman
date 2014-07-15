package scanner

import (
	"bufio"
)

func needMoreData() (int, []byte, error) { return 0, nil, nil }

func New(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		// TODO
		return 0, nil, nil
	})
	return scanner
}
