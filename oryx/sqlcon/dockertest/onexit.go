// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dockertest

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const interruptedExitCode = 130

// OnExit helps with cleaning up docker test.
type OnExit struct {
	sync.Mutex
	once     sync.Once
	handlers []func()
}

// NewOnExit create a new OnExit instance.
func NewOnExit() *OnExit {
	return &OnExit{
		handlers: make([]func(), 0),
	}
}

// Add adds a task that is executed on SIGINT, SIGKILL, SIGTERM.
func (at *OnExit) Add(f func()) {
	at.Lock()
	defer at.Unlock()
	at.handlers = append(at.handlers, f)
	at.once.Do(func() {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			<-c
			at.Exit(interruptedExitCode)
		}()
	})
}

// Exit wraps os.Exit
func (at *OnExit) Exit(status int) {
	at.execute()
	os.Exit(status)
}

func (at *OnExit) execute() {
	at.Lock()
	defer at.Unlock()
	for _, f := range at.handlers {
		f()
	}
	at.handlers = make([]func(), 0)
}
