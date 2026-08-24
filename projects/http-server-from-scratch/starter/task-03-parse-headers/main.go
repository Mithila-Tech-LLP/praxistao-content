package main

import (
	"bufio"
	"errors"
	"strings"
)

// ErrMalformedHeader is returned when a header line has no colon.
var ErrMalformedHeader = errors.New("malformed header")

// ParseHeaders reads HTTP headers from r until a blank line.
// Keys are lowercased. Returns ErrMalformedHeader on a bad line.
func ParseHeaders(r *bufio.Reader) (map[string]string, error) {
	headers := map[string]string{}
	// TODO:
	// Loop calling r.ReadString('\n') to read one line at a time.
	// Trim "\r\n" from the end of each line.
	// A blank line signals the end of headers — return headers, nil.
	// For non-blank lines, use strings.SplitN(line, ":", 2) to split.
	// If len(parts) != 2, return nil, ErrMalformedHeader.
	// Lowercase the key with strings.ToLower and trim spaces from key and value.
	_ = strings.ToLower
	_ = strings.TrimSpace
	return headers, nil
}

func main() {}
