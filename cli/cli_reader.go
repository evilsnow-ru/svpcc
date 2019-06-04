package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type StdioReader struct {
	helloPrefix string
	reader *bufio.Reader
}

func (r StdioReader) Read() ([]string, error) {
	fmt.Printf("%s => ", r.helloPrefix)
	data, err := r.reader.ReadBytes('\n')

	if err != nil {
		return nil, err
	} else {
		splitted := strings.Split(string(data), " ")
		idx := len(splitted) - 1
		splitted[idx] = strings.Trim(splitted[idx], "\n")
		return splitted, nil
	}
}

func NewReader(cmdPrefix string) *StdioReader {
	return &StdioReader{
		helloPrefix: cmdPrefix,
		reader: bufio.NewReader(os.Stdin),
	}
}
