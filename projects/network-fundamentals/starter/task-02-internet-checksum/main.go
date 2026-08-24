package main

// InternetChecksum computes the RFC 1071 Internet checksum of data.
func InternetChecksum(data []byte) uint16 {
	// TODO: sum every 16-bit big-endian word in data into a 32-bit accumulator.
	// TODO: if len(data) is odd, treat the trailing byte as the high byte of
	// a final word whose low byte is zero (pad right with a zero byte).
	// TODO: fold carries: while sum has bits above the low 16, add the
	// overflow back into the low 16 bits (sum = (sum & 0xFFFF) + (sum >> 16)).
	// TODO: return the one's complement (^uint16(sum)).
	return 0
}

func main() {}
