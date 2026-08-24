package main

import (
	"io"
	"net"
	"testing"
	"time"
)

func getFreeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()
	return addr
}

func TestEchoServer_EchoesData(t *testing.T) {
	addr := getFreeAddr(t)
	s := &EchoServer{Addr: addr}

	ln, err := s.ListenAndServe()
	if err != nil {
		t.Fatalf("ListenAndServe error: %v", err)
	}
	defer ln.Close()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))

	msg := "hello network\n"
	if _, err := conn.Write([]byte(msg)); err != nil {
		t.Fatalf("write error: %v", err)
	}

	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatalf("read error: %v", err)
	}
	if string(buf) != msg {
		t.Errorf("got %q, want %q", buf, msg)
	}
}

func TestEchoServer_MultipleWritesOnSameConnection(t *testing.T) {
	addr := getFreeAddr(t)
	s := &EchoServer{Addr: addr}

	ln, err := s.ListenAndServe()
	if err != nil {
		t.Fatalf("ListenAndServe error: %v", err)
	}
	defer ln.Close()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))

	for i := 0; i < 2; i++ {
		msg := "hello network\n"
		if _, err := conn.Write([]byte(msg)); err != nil {
			t.Fatalf("write %d error: %v", i, err)
		}
		buf := make([]byte, len(msg))
		if _, err := io.ReadFull(conn, buf); err != nil {
			t.Fatalf("read %d error: %v", i, err)
		}
		if string(buf) != msg {
			t.Errorf("iteration %d: got %q, want %q", i, buf, msg)
		}
	}
}
