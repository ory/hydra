// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package clock provides a small, injectable time source. Production code reads
// the current time through a Clock instead of calling time.Now() directly, so
// time-dependent behavior (such as flow expiry) can be tested deterministically.
package clock

import (
	"sync"
	"time"
)

// Clock is a source of the current time.
type Clock interface {
	Now() time.Time
}

// Provider is implemented by dependency containers (such as the Kratos registry)
// that expose a Clock. Embed it in a component's dependency interface to read
// the current time through the injected clock.
type Provider interface {
	Clock() Clock
}

// System is the production Clock backed by the real wall clock.
type System struct{}

// Now returns the current wall-clock time.
func (System) Now() time.Time { return time.Now() }

// New returns the production system clock.
func New() Clock { return System{} }

// Mock is a Clock whose time only changes when the caller advances it. It is
// safe for concurrent use: a test goroutine may Add while a request goroutine
// reads Now.
type Mock struct {
	mu sync.Mutex
	t  time.Time
}

// NewMock returns a Mock fixed at the given instant.
func NewMock(t time.Time) *Mock { return &Mock{t: t} }

// Now returns the Mock's current instant.
func (m *Mock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.t
}

// Add advances the Mock by d.
func (m *Mock) Add(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.t = m.t.Add(d)
}

// Set moves the Mock to the given instant.
func (m *Mock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.t = t
}
