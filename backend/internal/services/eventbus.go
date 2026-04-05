package services

import (
	"encoding/json"
	"sync"
)

// EventBus broadcasts events to SSE subscribers.
type EventBus struct {
	mu      sync.RWMutex
	clients map[chan Event]bool
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewEventBus() *EventBus {
	return &EventBus{
		clients: make(map[chan Event]bool),
	}
}

// Subscribe returns a channel that receives events. Call Unsubscribe to clean up.
func (eb *EventBus) Subscribe() chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, 16)
	eb.clients[ch] = true
	return ch
}

func (eb *EventBus) Unsubscribe(ch chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	delete(eb.clients, ch)
	close(ch)
}

// Publish sends an event to all subscribers.
func (eb *EventBus) Publish(eventType string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	event := Event{Type: eventType, Data: data}
	for ch := range eb.clients {
		select {
		case ch <- event:
		default:
			// channel full, skip (client too slow)
		}
	}
}

// PublishJSON publishes a raw JSON event.
func (eb *EventBus) PublishJSON(eventType string, data []byte) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	event := Event{Type: eventType}
	if err := json.Unmarshal(data, &event.Data); err != nil {
		event.Data = json.RawMessage(data)
	}
	for ch := range eb.clients {
		select {
		case ch <- event:
		default:
		}
	}
}

// ClientCount returns the number of active SSE subscribers.
func (eb *EventBus) ClientCount() int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.clients)
}
