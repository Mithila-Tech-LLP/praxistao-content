package main

import (
	"testing"
	"time"
)

func TestPublish_TwoSubscribers(t *testing.T) {
	ps := NewPubSub()
	ch1 := ps.Subscribe("news")
	ch2 := ps.Subscribe("news")

	count := ps.Publish("news", "hello")
	if count != 2 {
		t.Fatalf("expected count=2, got %d", count)
	}

	select {
	case msg := <-ch1:
		if msg != "hello" {
			t.Fatalf("ch1: expected hello, got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("ch1: timed out waiting for message")
	}

	select {
	case msg := <-ch2:
		if msg != "hello" {
			t.Fatalf("ch2: expected hello, got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("ch2: timed out waiting for message")
	}
}

func TestPublish_NoSubscribers(t *testing.T) {
	ps := NewPubSub()
	count := ps.Publish("empty", "msg")
	if count != 0 {
		t.Fatalf("expected count=0 for channel with no subscribers, got %d", count)
	}
}

func TestUnsubscribe(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("ch")
	ps.Unsubscribe("ch", ch)

	count := ps.Publish("ch", "after unsub")
	if count != 0 {
		t.Fatalf("expected count=0 after unsubscribe, got %d", count)
	}

	// No message should be buffered on the channel
	select {
	case msg := <-ch:
		t.Fatalf("expected no message after unsubscribe, got %q", msg)
	default:
		// ok — channel is empty
	}
}

func TestSubscribe_DifferentChannels(t *testing.T) {
	ps := NewPubSub()
	chA := ps.Subscribe("alpha")
	chB := ps.Subscribe("beta")

	ps.Publish("alpha", "for-alpha")

	select {
	case msg := <-chA:
		if msg != "for-alpha" {
			t.Fatalf("chA: expected for-alpha, got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("chA: timed out")
	}

	// chB should receive nothing
	select {
	case msg := <-chB:
		t.Fatalf("chB: should not receive message meant for alpha, got %q", msg)
	default:
		// ok
	}
}
