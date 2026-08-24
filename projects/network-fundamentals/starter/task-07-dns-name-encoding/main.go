package main

import "errors"

// ErrTruncatedDNSName is returned when data runs out before a terminating
// zero-length label is found.
var ErrTruncatedDNSName = errors.New("truncated dns name")

// EncodeDNSName encodes a dotted domain name into DNS wire format:
// length-prefixed labels terminated by a zero-length label.
func EncodeDNSName(name string) []byte {
	// TODO: split name on "." into labels.
	// TODO: for each label, append a single byte with its length, then the
	// label's raw bytes.
	// TODO: append a final zero byte to terminate the name.
	return nil
}

// DecodeDNSName decodes a wire-format DNS name starting at offset in data.
// It returns the reconstructed dotted name, the number of bytes consumed
// (including the terminating zero byte), and an error if data is truncated.
func DecodeDNSName(data []byte, offset int) (name string, bytesRead int, err error) {
	// TODO: starting at offset, repeatedly read a length byte, then that
	// many raw bytes as a label, until a zero length byte is read.
	// TODO: join the collected labels with "." to form the name.
	// TODO: bytesRead is the total number of bytes consumed, including the
	// final zero byte. Return ErrTruncatedDNSName if data runs out early.
	// Note: DNS compression pointers (0xC0 prefix) are out of scope here.
	return "", 0, errors.New("not implemented")
}

func main() {}
