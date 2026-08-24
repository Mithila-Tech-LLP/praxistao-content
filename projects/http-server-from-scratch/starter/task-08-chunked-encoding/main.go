package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// WriteChunked writes chunks using HTTP chunked transfer encoding.
// Each chunk: "<hex-length>\r\n<data>\r\n"
// Followed by the terminal chunk: "0\r\n\r\n"
func WriteChunked(w io.Writer, chunks [][]byte) error {
	// TODO:
	// For each chunk: fmt.Fprintf(w, "%x\r\n", len(chunk)), w.Write(chunk), fmt.Fprint(w, "\r\n")
	// Then write the terminal chunk: fmt.Fprint(w, "0\r\n\r\n")
	// Return any error encountered.
	_ = fmt.Fprintf
	return nil
}

// ReadChunked reads a chunked body until the terminal chunk.
// Returns all data concatenated.
func ReadChunked(r *bufio.Reader) ([]byte, error) {
	// TODO:
	// Loop:
	//   1. sizeLine, err := r.ReadString('\n'). Return err if non-nil.
	//   2. Parse hex size: strconv.ParseInt(strings.TrimSpace(sizeLine), 16, 64)
	//   3. If size == 0, read the trailing "\r\n" and return accumulated data.
	//   4. buf := make([]byte, size); io.ReadFull(r, buf)
	//   5. Read and discard the trailing "\r\n" after the data.
	//   6. Append buf to accumulated data.
	_ = strconv.ParseInt
	_ = strings.TrimSpace
	_ = io.ReadFull
	return nil, nil
}

func main() {}
