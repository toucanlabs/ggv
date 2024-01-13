package asty

import (
	"io"
	"os"
)

func noOpClose() error {
	return nil
}

func OpenRead(input string) (reader io.Reader, close func() error, err error) {
	if input == "" || input == "-" {
		return os.Stdin, noOpClose, nil
	}
	f, err := os.Open(input)
	if err != nil {
		return nil, nil, err
	}
	return f, func() error {
		return f.Close()
	}, nil
}
