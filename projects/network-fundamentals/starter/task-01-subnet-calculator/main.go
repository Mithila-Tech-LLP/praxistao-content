package main

import (
	"net"
)

// SubnetInfo parses cidr and returns the network address, broadcast address,
// first usable host, last usable host, and count of usable hosts.
func SubnetInfo(cidr string) (network, broadcast, firstHost, lastHost net.IP, hostCount int, err error) {
	// TODO: use net.ParseCIDR(cidr) to get the IP and *net.IPNet. Return its
	// error directly if parsing fails.
	// TODO: convert ipnet.IP to 4-byte form with .To4(), and treat the mask
	// and IP as 32-bit unsigned integers (encoding/binary.BigEndian is handy).
	// TODO: network = networkInt; broadcast = networkInt | ^maskInt.
	// TODO: firstHost = network + 1; lastHost = broadcast - 1.
	// TODO: hostCount = 2^(bits-ones) - 2, where ones, bits := mask.Size().
	return nil, nil, nil, nil, 0, nil
}

func main() {}
