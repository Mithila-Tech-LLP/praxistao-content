package main

import (
	"net"
	"strings"
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

func TestListenAndServe_RespondsOK(t *testing.T) {
	addr := getFreeAddr(t)
	s := &Server{Addr: addr}

	go s.ListenAndServe()
	time.Sleep(50 * time.Millisecond) // let the server start

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))
	conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"))

	buf := make([]byte, 512)
	n, _ := conn.Read(buf)
	resp := string(buf[:n])

	if !strings.HasPrefix(resp, "HTTP/1.1 200") {
		t.Errorf("got response %q, want prefix \"HTTP/1.1 200\"", resp)
	}
}
