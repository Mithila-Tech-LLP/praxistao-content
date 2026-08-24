package main

import "testing"

func TestSlidingWindow_FullStateMachine(t *testing.T) {
	w := NewSlidingWindow(4)

	for want := 0; want <= 3; want++ {
		seq, ok := w.Send()
		if !ok {
			t.Fatalf("Send() #%d: expected ok=true", want+1)
		}
		if seq != want {
			t.Errorf("Send() #%d: seq = %d, want %d", want+1, seq, want)
		}
	}

	if _, ok := w.Send(); ok {
		t.Error("Send() after 4 sends with no acks: expected ok=false (window full)")
	}

	w.Ack(1) // cumulative ack of seq 0 and 1; base becomes 2

	seq, ok := w.Send()
	if !ok || seq != 4 {
		t.Errorf("Send() after Ack(1) = (%d, %v), want (4, true)", seq, ok)
	}

	seq, ok = w.Send()
	if !ok || seq != 5 {
		t.Errorf("second Send() after Ack(1) = (%d, %v), want (5, true)", seq, ok)
	}

	if _, ok := w.Send(); ok {
		t.Error("Send() after window refilled to base=2,next=6: expected ok=false")
	}
}

func TestSlidingWindow_StaleAckIgnored(t *testing.T) {
	w := NewSlidingWindow(4)
	w.Send()
	w.Send()
	w.Ack(1) // base -> 2
	w.Ack(0) // stale, should be a no-op (would move base backward)

	// base should still be 2: only 2 more sends should be possible before
	// hitting the window limit relative to base=2.
	seq, ok := w.Send()
	if !ok || seq != 2 {
		t.Errorf("Send() = (%d, %v), want (2, true)", seq, ok)
	}
}
