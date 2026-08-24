package main

import (
	"errors"
	"net"
)

// EchoServer is a TCP server that echoes back everything it reads.
type EchoServer struct{ Addr string }

// ListenAndServe opens a TCP listener on s.Addr and starts accepting
// connections in the background, returning the listener immediately.
func (s *EchoServer) ListenAndServe() (net.Listener, error) {
	// TODO: open a listener with net.Listen("tcp", s.Addr).
	// TODO: in a goroutine, loop calling ln.Accept(); for each accepted
	// connection, handle it in its own goroutine.
	// TODO: in the per-connection goroutine, copy everything read from the
	// connection back to the connection (io.Copy(conn, conn) works) until
	// the client closes or an error occurs, then close the connection.
	// TODO: return the listener (not nil) immediately so the caller can
	// close it to stop the server.
	return nil, errors.New("not implemented")
}

func main() {}
