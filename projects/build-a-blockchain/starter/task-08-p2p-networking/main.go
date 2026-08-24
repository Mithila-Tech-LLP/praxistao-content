package main

import (
	"net"
)

type MessageType byte

const (
	MsgBlock MessageType = iota
	MsgTx
)

// SendMessage writes a length-prefixed frame to conn: one byte for
// msgType, then a 4-byte big-endian length, then payload.
func SendMessage(conn net.Conn, msgType MessageType, payload []byte) error {
	panic("TODO: implement SendMessage")
}

// ReadMessage reads one complete length-prefixed frame from conn and
// returns its type and payload.
func ReadMessage(conn net.Conn) (MessageType, []byte, error) {
	panic("TODO: implement ReadMessage using io.ReadFull")
}

// Node listens on an address and invokes onMessage for every message it
// receives on any connection, in a new goroutine per connection.
type Node struct {
	Address   string
	onMessage func(MessageType, []byte)
	listener  net.Listener
}

func NewNode(address string, onMessage func(MessageType, []byte)) *Node {
	return &Node{Address: address, onMessage: onMessage}
}

// Listen starts accepting connections in the background. It returns
// once the listener is bound (or an error if binding failed). After a
// successful call, n.Address holds the real, resolved address (useful
// when the caller passed a ":0" port).
func (n *Node) Listen() error {
	panic("TODO: implement Listen -- net.Listen, then start an accept loop goroutine")
}

// Send connects to peerAddress and sends one message.
func (n *Node) Send(peerAddress string, msgType MessageType, payload []byte) error {
	panic("TODO: implement Send")
}

// Close stops accepting new connections.
func (n *Node) Close() error {
	panic("TODO: implement Close")
}
