package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// EventType represents the type of data change event
type EventType string

const (
	EventTypeTreeChanged   EventType = "tree_changed"
	EventTypePromptChanged EventType = "prompt_changed"
	EventTypeNodeChanged   EventType = "node_changed"
	EventTypeNoteChanged   EventType = "note_changed"
)

// Event represents a notification event
type Event struct {
	Type      EventType `json:"type"`
	PromptID  *int      `json:"prompt_id,omitempty"`
	Message   string    `json:"message"`
	Timestamp int64     `json:"timestamp"`
}

// Client represents a connected SSE client
type Client struct {
	ID       string
	Send     chan Event
	Done     chan bool
	IsActive bool
}

// Notifier manages SSE connections and broadcasts events
type Notifier struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

// NewNotifier creates a new notification broadcaster
func NewNotifier() *Notifier {
	return &Notifier{
		clients: make(map[string]*Client),
	}
}

// RegisterClient adds a new client to receive notifications
func (n *Notifier) RegisterClient(clientID string) *Client {
	n.mu.Lock()
	defer n.mu.Unlock()

	client := &Client{
		ID:       clientID,
		Send:     make(chan Event, 256), // Buffered channel
		Done:     make(chan bool),
		IsActive: true,
	}

	n.clients[clientID] = client
	return client
}

// UnregisterClient removes a client from receiving notifications
func (n *Notifier) UnregisterClient(clientID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if client, exists := n.clients[clientID]; exists {
		close(client.Send)
		delete(n.clients, clientID)
	}
}

// Broadcast sends an event to all connected clients
func (n *Notifier) Broadcast(event Event) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Send event to all active clients
	for _, client := range n.clients {
		if client.IsActive {
			select {
			case client.Send <- event:
				// Event sent successfully
			default:
				// Channel is full, client might be slow
				// Mark as inactive but don't remove yet
				client.IsActive = false
			}
		}
	}
}

// BroadcastTreeChanged notifies all clients that the tree structure changed
func (n *Notifier) BroadcastTreeChanged() {
	event := Event{
		Type:      EventTypeTreeChanged,
		Message:   "Tree structure has been updated",
		Timestamp: time.Now().Unix(),
	}
	n.Broadcast(event)
}

// BroadcastPromptChanged notifies all clients that a prompt changed
func (n *Notifier) BroadcastPromptChanged(promptID int) {
	event := Event{
		Type:      EventTypePromptChanged,
		PromptID:  &promptID,
		Message:   fmt.Sprintf("Prompt %d has been updated", promptID),
		Timestamp: time.Now().Unix(),
	}
	n.Broadcast(event)
}

// BroadcastNodeChanged notifies all clients that a node changed
func (n *Notifier) BroadcastNodeChanged(promptID int) {
	event := Event{
		Type:      EventTypeNodeChanged,
		PromptID:  &promptID,
		Message:   fmt.Sprintf("Nodes for prompt %d have been updated", promptID),
		Timestamp: time.Now().Unix(),
	}
	n.Broadcast(event)
}

// BroadcastNoteChanged notifies all clients that a note changed
func (n *Notifier) BroadcastNoteChanged(promptID int) {
	event := Event{
		Type:      EventTypeNoteChanged,
		PromptID:  &promptID,
		Message:   fmt.Sprintf("Notes for prompt %d have been updated", promptID),
		Timestamp: time.Now().Unix(),
	}
	n.Broadcast(event)
}

// GetActiveClientsCount returns the number of active clients
func (n *Notifier) GetActiveClientsCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.clients)
}

// FormatSSE formats an event as Server-Sent Events format
func FormatSSE(event Event) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data: %s\n\n", string(data)), nil
}

