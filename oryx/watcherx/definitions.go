// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type (
	errSchemeUnknown struct {
		scheme string
	}
	EventChannel chan Event
	Watcher      interface {
		// DispatchNow fires the watcher and causes an event.
		//
		// WARNING: The returned channel must be read or no further events will
		// be propagated due to a deadlock.
		DispatchNow() (<-chan int, error)
	}
	dispatcher struct {
		ctx     context.Context
		trigger chan struct{}
		done    chan int
	}
)

var (
	// ErrSchemeUnknown is just for checking with errors.Is()
	ErrSchemeUnknown     = &errSchemeUnknown{}
	ErrWatcherNotRunning = fmt.Errorf("watcher is not running")
)

func (e *errSchemeUnknown) Is(other error) bool {
	_, ok := other.(*errSchemeUnknown)
	return ok
}

func (e *errSchemeUnknown) Error() string {
	return fmt.Sprintf("unknown scheme '%s' to watch", e.scheme)
}

func newDispatcher(ctx context.Context) *dispatcher {
	return &dispatcher{
		ctx:     ctx,
		trigger: make(chan struct{}),
		done:    make(chan int),
	}
}

func (d *dispatcher) DispatchNow() (<-chan int, error) {
	if d.trigger == nil {
		return nil, ErrWatcherNotRunning
	}
	// The trigger send must respect cancellation. Once the watcher's context is
	// done its receiver goroutine has returned, so a bare send would block
	// forever and wedge the caller.
	select {
	case d.trigger <- struct{}{}:
		return d.done, nil
	case <-d.ctx.Done():
		return nil, errors.WithStack(ErrWatcherNotRunning)
	}
}

func Watch(ctx context.Context, u *url.URL, c EventChannel) (Watcher, error) {
	switch u.Scheme {
	// see urlx.Parse for why the empty string is also file
	case "file", "":
		return WatchFile(ctx, u.Path, c)
	}
	return nil, &errSchemeUnknown{u.Scheme}
}
