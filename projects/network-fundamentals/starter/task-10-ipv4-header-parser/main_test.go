package main

import (
	"net"
	"testing"
)

func TestParseIPv4Header_Valid(t *testing.T) {
	data := []byte{0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0x00, 0x00, 0xac, 0x10, 0x0a, 0x63, 0xac, 0x10, 0x0a, 0x0c}

	version, ihl, totalLength, ttl, protocol, srcIP, dstIP, err := ParseIPv4Header(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != 4 {
		t.Errorf("version = %d, want 4", version)
	}
	if ihl != 5 {
		t.Errorf("ihl = %d, want 5", ihl)
	}
	if totalLength != 60 {
		t.Errorf("totalLength = %d, want 60", totalLength)
	}
	if ttl != 64 {
		t.Errorf("ttl = %d, want 64", ttl)
	}
	if protocol != 6 {
		t.Errorf("protocol = %d, want 6", protocol)
	}
	if !srcIP.Equal(net.ParseIP("172.16.10.99")) {
		t.Errorf("srcIP = %v, want 172.16.10.99", srcIP)
	}
	if !dstIP.Equal(net.ParseIP("172.16.10.12")) {
		t.Errorf("dstIP = %v, want 172.16.10.12", dstIP)
	}
}

func TestParseIPv4Header_TooShort(t *testing.T) {
	data := make([]byte, 10)
	_, _, _, _, _, _, _, err := ParseIPv4Header(data)
	if err == nil {
		t.Error("expected error for header shorter than 20 bytes")
	}
}
