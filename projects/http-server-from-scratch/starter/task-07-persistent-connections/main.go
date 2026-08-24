package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Response is the HTTP response returned by handler functions.
type Response struct {
	Status int
	Body   string
}

// ServeConn handles persistent HTTP/1.1 connections.
// It reads requests in a loop, dispatches each to the dispatch function,
// and writes the response back. It stops when:
//   - The request includes "Connection: close"
//   - Reading fails (EOF or error)
func ServeConn(conn net.Conn, dispatch func(method, path string, headers map[string]string) Response) {
	// TODO:
	// 1. reader := bufio.NewReader(conn)
	// 2. Loop:
	//    a. Read the request line with reader.ReadString('\n'). Break on error.
	//    b. Parse method/path from the request line (strings.Fields).
	//    c. Read headers: loop reader.ReadString('\n') until blank line,
	//       build headers map (lowercase keys, split on ":").
	//    d. If headers["connection"] == "close", write response then break.
	//    e. Call dispatch(method, path, headers).
	//    f. Write "HTTP/1.1 <status> OK\r\n\r\n<body>" to conn.
	_ = bufio.NewReader(conn)
	_ = strings.ToLower
	_ = fmt.Fprintf
}

func main() {}
