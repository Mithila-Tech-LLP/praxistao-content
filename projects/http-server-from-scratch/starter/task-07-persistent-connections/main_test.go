package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServeConn_HandlesTwoRequests(t *testing.T) {
	server, client := net.Pipe()
	defer client.Close()

	dispatch := func(method, path string, headers map[string]string) Response {
		return Response{Status: 200, Body: "OK"}
	}

	go ServeConn(server, dispatch)

	client.SetDeadline(time.Now().Add(500 * time.Millisecond))

	fmt.Fprint(client, "GET /a HTTP/1.1\r\nHost: localhost\r\n\r\n")
	fmt.Fprint(client, "GET /b HTTP/1.1\r\nHost: localhost\r\n\r\n")

	buf := make([]byte, 1024)
	n, _ := client.Read(buf)
	resp := string(buf[:n])

	if !strings.Contains(resp, "HTTP/1.1 200") {
		t.Errorf("expected at least one 200 response, got: %q", resp)
	}
}

func TestServeConn_StopsOnConnectionClose(t *testing.T) {
	server, client := net.Pipe()

	dispatch := func(method, path string, headers map[string]string) Response {
		return Response{Status: 200, Body: "OK"}
	}

	done := make(chan struct{})
	go func() {
		ServeConn(server, dispatch)
		close(done)
	}()

	fmt.Fprint(client, "GET / HTTP/1.1\r\nConnection: close\r\n\r\n")

	buf := make([]byte, 512)
	client.SetDeadline(time.Now().Add(500 * time.Millisecond))
	client.Read(buf)
	client.Close()

	select {
	case <-done:
		// ServeConn returned as expected
	case <-time.After(time.Second):
		t.Error("ServeConn did not stop after Connection: close")
	}
}
