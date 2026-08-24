package main

// SlidingWindow models a TCP-style sliding send window with cumulative ACKs.
type SlidingWindow struct {
	base int // oldest unacknowledged sequence number
	next int // next sequence number to send
	size int // window size
}

// NewSlidingWindow returns a window of the given size, starting at seq 0.
func NewSlidingWindow(size int) *SlidingWindow {
	// TODO: return &SlidingWindow{base: 0, next: 0, size: size}.
	return &SlidingWindow{size: size}
}

// Send returns the next sequence number to send, or ok=false if the window
// is full (no room to send until an ACK arrives).
func (w *SlidingWindow) Send() (seq int, ok bool) {
	// TODO: only allow sending while w.next < w.base+w.size.
	// TODO: on success, return the current w.next and then increment it.
	return 0, false
}

// Ack cumulatively acknowledges everything up to and including seq.
func (w *SlidingWindow) Ack(seq int) {
	// TODO: advance w.base to seq+1, but only if that's forward progress —
	// ignore stale/out-of-order acks that would move base backward.
}

func main() {}
