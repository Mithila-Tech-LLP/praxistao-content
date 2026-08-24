package main

import "net"

// Server holds the configuration for our TCP server.
type Server struct{ Addr string }

// ListenAndServe opens a TCP listener on s.Addr, accepts one connection,
// reads up to 1024 bytes, writes a fixed HTTP 200 response, and closes.
func (s *Server) ListenAndServe() error {
	// TODO:
	// 1. net.Listen("tcp", s.Addr)
	// 2. ln.Accept() to get a net.Conn
	// 3. Read up to 1024 bytes with conn.Read(buf) — ignore the error
	// 4. Write "HTTP/1.1 200 OK\r\n\r\nOK\n" to conn
	// 5. conn.Close(), return nil
	_ = net.Listen
	return nil
}

func main() {}
