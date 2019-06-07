package cli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
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
		splitted := strings.Split(strings.TrimSpace(string(data)), " ")
		idx := len(splitted) - 1

		var cutstring string

		if runtime.GOOS == "windows" {
			cutstring = "\r\n"
		} else {
			cutstring = "\n"
		}

		splitted[idx] = strings.Trim(splitted[idx], cutstring)
		return splitted, nil
	}
}

func NewReader(cmdPrefix string) *StdioReader {
	return &StdioReader{
		helloPrefix: cmdPrefix,
		reader: bufio.NewReader(os.Stdin),
	}
}
