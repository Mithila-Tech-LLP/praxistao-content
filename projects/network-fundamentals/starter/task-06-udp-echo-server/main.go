package main

import (
	"errors"
	"net"
)

// StartUDPEchoServer opens a UDP socket on addr and echoes every datagram
// it receives back to its sender, returning the connection immediately.
func StartUDPEchoServer(addr string) (*net.UDPConn, error) {
	// TODO: resolve addr with net.ResolveUDPAddr("udp", addr) and open the
	// socket with net.ListenUDP("udp", udpAddr).
	// TODO: in a goroutine, loop: ReadFromUDP into a buffer, then
	// WriteToUDP the same bytes back to the sender's address. Repeat
	// forever (don't stop after the first packet).
	// TODO: return the *net.UDPConn immediately so the caller can Close()
	// it to stop the server.
	return nil, errors.New("not implemented")
}

func main() {}
