// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package reqlog

import (
	"context"
	"sync"
	"time"
)

// WithEnableExternalLatencyMeasurement returns a context that measures external latencies.
func WithEnableExternalLatencyMeasurement(ctx context.Context) context.Context {
	container := contextContainer{
		latencies: make([]externalLatency, 0),
	}
	return context.WithValue(ctx, externalLatencyKey, &container)
}

// StartMeasureExternalCall starts measuring the duration of an external call.
// The returned function has to be called to record the duration.
func StartMeasureExternalCall(ctx context.Context, cause, detail string, start time.Time) {
	container, ok := ctx.Value(externalLatencyKey).(*contextContainer)
	if !ok {
		return
	}
	if _, ok := ctx.Value(disableExternalLatencyMeasurement).(bool); ok {
		return
	}

	container.Lock()
	defer container.Unlock()
	container.latencies = append(container.latencies, externalLatency{
		Took:   time.Since(start),
		Cause:  cause,
		Detail: detail,
	})
}

// totalExternalLatency returns the total duration of all external calls.
func totalExternalLatency(ctx context.Context) (total time.Duration) {
	if _, ok := ctx.Value(disableExternalLatencyMeasurement).(bool); ok {
		return 0
	}
	container, ok := ctx.Value(externalLatencyKey).(*contextContainer)
	if !ok {
		return 0
	}

	container.Lock()
	defer container.Unlock()
	for _, l := range container.latencies {
		total += l.Took
	}
	return total
}

// WithDisableExternalLatencyMeasurement returns a context that does not measure external latencies.
// Use this when you want to disable external latency measurements for a specific request.
func WithDisableExternalLatencyMeasurement(ctx context.Context) context.Context {
	return context.WithValue(ctx, disableExternalLatencyMeasurement, true)
}

type (
	externalLatency = struct {
		Took          time.Duration
		Cause, Detail string
	}
	contextContainer = struct {
		latencies []externalLatency
		sync.Mutex
	}
	contextKey int
)

const (
	externalLatencyKey                contextKey = 1
	disableExternalLatencyMeasurement contextKey = 2
)
