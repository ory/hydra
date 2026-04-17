// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package mailhog is a copy of github.com/mailhog/storage with the missing
// read mutex operations added to List, Count, Search, and Load.
package mailhog

import (
	"errors"
	"strings"
	"sync"

	"github.com/mailhog/data"
)

// InMemory is an in memory storage backend
type InMemory struct {
	MessageIDIndex map[string]int
	Messages       []*data.Message
	mu             sync.RWMutex
}

// NewInMemory creates a new in memory storage backend
func NewInMemory() *InMemory {
	return &InMemory{
		MessageIDIndex: make(map[string]int),
		Messages:       make([]*data.Message, 0),
	}
}

// Store stores a message and returns its storage ID
func (memory *InMemory) Store(m *data.Message) (string, error) {
	memory.mu.Lock()
	defer memory.mu.Unlock()
	memory.Messages = append(memory.Messages, m)
	memory.MessageIDIndex[string(m.ID)] = len(memory.Messages) - 1
	return string(m.ID), nil
}

// Count returns the number of stored messages
func (memory *InMemory) Count() int {
	memory.mu.RLock()
	defer memory.mu.RUnlock()
	return len(memory.Messages)
}

// Search finds messages matching the query
func (memory *InMemory) Search(kind, query string, start, limit int) (*data.Messages, int, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()
	// FIXME needs optimising, or replacing with a proper db!
	query = strings.ToLower(query)
	var filteredMessages = make([]*data.Message, 0)
	for _, m := range memory.Messages {
		doAppend := false

		switch kind {
		case "to":
			for _, to := range m.To {
				if strings.Contains(strings.ToLower(to.Mailbox+"@"+to.Domain), query) {
					doAppend = true
					break
				}
			}
			if !doAppend {
				if hdr, ok := m.Content.Headers["To"]; ok {
					for _, to := range hdr {
						if strings.Contains(strings.ToLower(to), query) {
							doAppend = true
							break
						}
					}
				}
			}
		case "from":
			if strings.Contains(strings.ToLower(m.From.Mailbox+"@"+m.From.Domain), query) {
				doAppend = true
			}
			if !doAppend {
				if hdr, ok := m.Content.Headers["From"]; ok {
					for _, from := range hdr {
						if strings.Contains(strings.ToLower(from), query) {
							doAppend = true
							break
						}
					}
				}
			}
		case "containing":
			if strings.Contains(strings.ToLower(m.Content.Body), query) {
				doAppend = true
			}
			if !doAppend {
				for _, hdr := range m.Content.Headers {
					for _, v := range hdr {
						if strings.Contains(strings.ToLower(v), query) {
							doAppend = true
						}
					}
				}
			}
		}

		if doAppend {
			filteredMessages = append(filteredMessages, m)
		}
	}

	var messages = make([]data.Message, 0)

	if len(filteredMessages) == 0 || start > len(filteredMessages) {
		msgs := data.Messages(messages)
		return &msgs, 0, nil
	}

	if start+limit > len(filteredMessages) {
		limit = len(filteredMessages) - start
	}

	start = len(filteredMessages) - start - 1
	end := start - limit

	if start < 0 {
		start = 0
	}
	if end < -1 {
		end = -1
	}

	for i := start; i > end; i-- {
		//for _, m := range memory.MessageIndex[start:end] {
		messages = append(messages, *filteredMessages[i])
	}

	msgs := data.Messages(messages)
	return &msgs, len(filteredMessages), nil
}

// List lists stored messages by index
func (memory *InMemory) List(start int, limit int) (*data.Messages, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()
	var messages = make([]data.Message, 0)

	if len(memory.Messages) == 0 || start > len(memory.Messages) {
		msgs := data.Messages(messages)
		return &msgs, nil
	}

	if start+limit > len(memory.Messages) {
		limit = len(memory.Messages) - start
	}

	start = len(memory.Messages) - start - 1
	end := start - limit

	if start < 0 {
		start = 0
	}
	if end < -1 {
		end = -1
	}

	for i := start; i > end; i-- {
		//for _, m := range memory.MessageIndex[start:end] {
		messages = append(messages, *memory.Messages[i])
	}

	msgs := data.Messages(messages)
	return &msgs, nil
}

// DeleteOne deletes an individual message by storage ID
func (memory *InMemory) DeleteOne(id string) error {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	var index int
	var ok bool

	if index, ok = memory.MessageIDIndex[id]; !ok && true {
		return errors.New("message not found")
	}

	delete(memory.MessageIDIndex, id)
	for k, v := range memory.MessageIDIndex {
		if v > index {
			memory.MessageIDIndex[k] = v - 1
		}
	}
	memory.Messages = append(memory.Messages[:index], memory.Messages[index+1:]...)
	return nil
}

// DeleteAll deletes all in memory messages
func (memory *InMemory) DeleteAll() error {
	memory.mu.Lock()
	defer memory.mu.Unlock()
	memory.Messages = make([]*data.Message, 0)
	memory.MessageIDIndex = make(map[string]int)
	return nil
}

// Load returns an individual message by storage ID
func (memory *InMemory) Load(id string) (*data.Message, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()
	if idx, ok := memory.MessageIDIndex[id]; ok {
		return memory.Messages[idx], nil
	}
	return nil, nil
}
