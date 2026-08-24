package main

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestSendAndReadMessageRoundTrip(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	done := make(chan struct{})
	var gotType MessageType
	var gotPayload []byte

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		gotType, gotPayload, _ = ReadMessage(conn)
		close(done)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	if err := SendMessage(conn, MsgTx, []byte("hello peer")); err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for message to be read")
	}

	if gotType != MsgTx {
		t.Fatalf("expected MsgTx, got %v", gotType)
	}
	if !bytes.Equal(gotPayload, []byte("hello peer")) {
		t.Fatalf("expected payload %q, got %q", "hello peer", gotPayload)
	}
}

func TestNodeReceivesMessage(t *testing.T) {
	received := make(chan struct {
		Type    MessageType
		Payload []byte
	}, 1)

	node := NewNode("127.0.0.1:0", func(msgType MessageType, payload []byte) {
		received <- struct {
			Type    MessageType
			Payload []byte
		}{msgType, payload}
	})

	if err := node.Listen(); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer node.Close()

	conn, err := net.Dial("tcp", node.Address)
	if err != nil {
		t.Fatalf("failed to dial node: %v", err)
	}
	defer conn.Close()

	if err := SendMessage(conn, MsgBlock, []byte("a block")); err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	select {
	case msg := <-received:
		if msg.Type != MsgBlock {
			t.Fatalf("expected MsgBlock, got %v", msg.Type)
		}
		if !bytes.Equal(msg.Payload, []byte("a block")) {
			t.Fatalf("expected payload %q, got %q", "a block", msg.Payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for node to receive the message")
	}
}

func TestTwoNodesCommunicate(t *testing.T) {
	receivedB := make(chan []byte, 1)

	nodeA := NewNode("127.0.0.1:0", nil)
	nodeB := NewNode("127.0.0.1:0", func(msgType MessageType, payload []byte) {
		receivedB <- payload
	})

	if err := nodeA.Listen(); err != nil {
		t.Fatalf("nodeA Listen failed: %v", err)
	}
	defer nodeA.Close()

	if err := nodeB.Listen(); err != nil {
		t.Fatalf("nodeB Listen failed: %v", err)
	}
	defer nodeB.Close()

	if err := nodeA.Send(nodeB.Address, MsgTx, []byte("tx from A")); err != nil {
		t.Fatalf("nodeA.Send failed: %v", err)
	}

	select {
	case payload := <-receivedB:
		if !bytes.Equal(payload, []byte("tx from A")) {
			t.Fatalf("expected %q, got %q", "tx from A", payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for nodeB to receive nodeA's message")
	}
}

func TestEmptyPayload(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, payload, err := ReadMessage(conn)
		if err != nil {
			t.Errorf("ReadMessage failed: %v", err)
		}
		if len(payload) != 0 {
			t.Errorf("expected an empty payload, got %d bytes", len(payload))
		}
		close(done)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	if err := SendMessage(conn, MsgBlock, nil); err != nil {
		t.Fatalf("SendMessage with nil payload failed: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out")
	}
}
