# Task 07: Pub/Sub

## What you will build

A publish/subscribe message bus. Subscribers listen on named channels; publishers fan out messages to every subscriber on a channel without knowing who they are. This is how Redis Pub/Sub works and how event-driven architectures decouple producers from consumers.

## Concepts

### Channels as Go channels

Each subscriber gets a `chan string`. When `Publish` is called, it sends the message on every subscriber's channel. The subscriber reads from its channel in its own goroutine.

```
Publish("news", "hello")
    → subscriber A receives "hello"
    → subscriber B receives "hello"
```

### Fan-out

Store subscribers as `map[string][]chan string` — channel name → list of subscriber channels. `Publish` iterates over the list and sends to each.

### Non-blocking sends

If a subscriber's goroutine is slow, a blocking send would stall `Publish` for all other subscribers. Use a buffered channel (e.g. buffer 16) so `Publish` does not block for a momentarily busy consumer.

## Interface to implement

```go
// Subscribe creates a new subscription to channel.
// Returns a receive-only channel. The caller reads messages from it.
func (s *Store) Subscribe(channel string) <-chan string

// Publish sends message to all subscribers of channel.
// Returns the number of subscribers that received the message.
func (s *Store) Publish(channel, message string) int

// Unsubscribe removes the subscription represented by ch from channel.
// The channel is closed after removal.
func (s *Store) Unsubscribe(channel string, ch <-chan string)
```

## Hints

- Use a separate mutex for pub/sub state (or a dedicated `sync.RWMutex`) rather than sharing the key/value mutex. Pub/sub operations should not block key reads.
- After `Unsubscribe`, close the removed channel so the subscriber's `range` loop can terminate.
- Do not send on a closed channel — that panics. Mark the channel as unsubscribed before closing it.

## Run the tests

```bash
cd starter/task-07-pub-sub
go test ./...
```
