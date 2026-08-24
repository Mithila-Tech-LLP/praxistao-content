package main

import (
	"errors"
	"net"
)

// ErrHeaderTooShort is returned when data is shorter than a minimal
// (no-options) 20-byte IPv4 header.
var ErrHeaderTooShort = errors.New("ipv4 header too short")

// ParseIPv4Header parses a raw IPv4 header with no options (exactly 20
// bytes) out of data.
func ParseIPv4Header(data []byte) (version, ihl, totalLength, ttl, protocol int, srcIP, dstIP net.IP, err error) {
	// TODO: return ErrHeaderTooShort if len(data) < 20.
	// TODO: version = high nibble of data[0]; ihl = low nibble of data[0].
	// TODO: totalLength = big-endian uint16 at data[2:4].
	// TODO: ttl = data[8]; protocol = data[9].
	// TODO: srcIP = data[12:16] as a net.IP; dstIP = data[16:20] as a net.IP.
	return 0, 0, 0, 0, 0, nil, nil, errors.New("not implemented")
}

func main() {}
