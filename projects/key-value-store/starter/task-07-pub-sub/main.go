package main

import "sync"

type PubSub struct {
	mu   sync.Mutex
	subs map[string][]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{subs: make(map[string][]chan string)}
}

func (ps *PubSub) Subscribe(channel string) <-chan string {
	// TODO: create buffered channel (size 10), register it, return it
	return nil
}

func (ps *PubSub) Publish(channel, message string) int {
	// TODO: send message to all subscribers, return count
	return 0
}

func (ps *PubSub) Unsubscribe(channel string, ch <-chan string) {
	// TODO: remove ch from channel's subscriber list
}
