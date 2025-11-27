// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package reqlog

import (
	"context"
	"sync/atomic"
	"time"
)

func withEnableExternalLatencyMeasurement(ctx context.Context) context.Context {
	return context.WithValue(ctx, externalLatencyKey, new(int64))
}

func AccumulateExternalLatency(ctx context.Context, dur time.Duration) {
	total, ok := ctx.Value(externalLatencyKey).(*int64)
	if !ok {
		return
	}
	atomic.AddInt64(total, int64(dur))
}

func getExternalLatency(ctx context.Context) time.Duration {
	total, ok := ctx.Value(externalLatencyKey).(*int64)
	if !ok {
		return 0
	}
	return time.Duration(atomic.LoadInt64(total))
}

type contextKey int

const externalLatencyKey contextKey = 1
