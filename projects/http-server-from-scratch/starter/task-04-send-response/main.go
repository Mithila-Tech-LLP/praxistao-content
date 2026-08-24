package main

import (
	"fmt"
	"io"
)

// Response represents an HTTP response.
type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

// statusText maps common status codes to their reason phrase.
func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}

// Write serializes the HTTP/1.1 response to w.
func (r *Response) Write(w io.Writer) error {
	// TODO:
	// 1. fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", r.Status, statusText(r.Status))
	// 2. For each key/value in r.Headers: fmt.Fprintf(w, "%s: %s\r\n", key, value)
	// 3. fmt.Fprint(w, "\r\n")
	// 4. w.Write(r.Body)
	// Return any error from the writes above.
	_ = fmt.Fprintf
	return nil
}

func main() {}
