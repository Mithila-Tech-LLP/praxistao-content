package main

import "errors"

// ErrMalformedRequestLine is returned when the request line cannot be parsed.
var ErrMalformedRequestLine = errors.New("malformed request line")

// ParseRequestLine parses an HTTP request line like "GET /path HTTP/1.1".
// Returns method, path, version or ErrMalformedRequestLine.
func ParseRequestLine(line string) (method, path, version string, err error) {
	// TODO: use strings.Fields to split line into parts.
	// Return ErrMalformedRequestLine if len(parts) != 3.
	// Otherwise return parts[0], parts[1], parts[2], nil.
	return "", "", "", ErrMalformedRequestLine
}

func main() {}
