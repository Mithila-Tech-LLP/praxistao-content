package main

import (
	"bufio"
	"io"
	"strconv"
)

// ParseBody reads exactly Content-Length bytes from r.
// Returns empty slice if Content-Length is absent or zero.
func ParseBody(r *bufio.Reader, headers map[string]string) ([]byte, error) {
	// TODO:
	// 1. Look up headers["content-length"]. If absent, return []byte{}, nil.
	// 2. strconv.Atoi to parse the length. Return error if invalid.
	// 3. If length is 0, return []byte{}, nil.
	// 4. buf := make([]byte, length); io.ReadFull(r, buf)
	// 5. Return buf, err.
	_ = io.ReadFull
	_ = strconv.Atoi
	return []byte{}, nil
}

func main() {}
