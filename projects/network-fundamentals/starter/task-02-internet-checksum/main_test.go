package main

import "testing"

func TestInternetChecksum_IPv4Header(t *testing.T) {
	data := []byte{0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0x00, 0x00, 0xac, 0x10, 0x0a, 0x63, 0xac, 0x10, 0x0a, 0x0c}
	got := InternetChecksum(data)
	if got != 0xb1e6 {
		t.Errorf("got 0x%x, want 0xb1e6", got)
	}
}

func TestInternetChecksum_AllZero(t *testing.T) {
	got := InternetChecksum([]byte{0, 0, 0, 0})
	if got != 0xffff {
		t.Errorf("got 0x%x, want 0xffff", got)
	}
}

func TestInternetChecksum_OddLength(t *testing.T) {
	got := InternetChecksum([]byte{0x00, 0x01, 0x02})
	if got != 0xfdfe {
		t.Errorf("got 0x%x, want 0xfdfe", got)
	}
}

func TestInternetChecksum_DeadBeef(t *testing.T) {
	got := InternetChecksum([]byte{0xde, 0xad, 0xbe, 0xef})
	if got != 0x6262 {
		t.Errorf("got 0x%x, want 0x6262", got)
	}
}
