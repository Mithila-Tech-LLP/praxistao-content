package main

import (
	"net"
	"testing"
	"time"
)

func TestUDPEchoServer_EchoesDatagram(t *testing.T) {
	server, err := StartUDPEchoServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("StartUDPEchoServer error: %v", err)
	}
	defer server.Close()

	serverAddr := server.LocalAddr().(*net.UDPAddr)

	client, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		t.Fatalf("DialUDP error: %v", err)
	}
	defer client.Close()

	client.SetDeadline(time.Now().Add(time.Second))

	msg := []byte("ping")
	if _, err := client.Write(msg); err != nil {
		t.Fatalf("write error: %v", err)
	}

	buf := make([]byte, 512)
	n, err := client.Read(buf)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if string(buf[:n]) != "ping" {
		t.Errorf("got %q, want %q", buf[:n], "ping")
	}
}

func TestUDPEchoServer_TwoDatagrams(t *testing.T) {
	server, err := StartUDPEchoServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("StartUDPEchoServer error: %v", err)
	}
	defer server.Close()

	serverAddr := server.LocalAddr().(*net.UDPAddr)

	client, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		t.Fatalf("DialUDP error: %v", err)
	}
	defer client.Close()

	client.SetDeadline(time.Now().Add(time.Second))

	messages := []string{"first", "second"}
	for _, msg := range messages {
		if _, err := client.Write([]byte(msg)); err != nil {
			t.Fatalf("write error: %v", err)
		}
		buf := make([]byte, 512)
		n, err := client.Read(buf)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}
		if string(buf[:n]) != msg {
			t.Errorf("got %q, want %q", buf[:n], msg)
		}
	}
}
