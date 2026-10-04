// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
)

// ExemplarFromContext returns exemplar labels (trace_id, span_id) for the
// sampled trace span in ctx, or nil when there is none. The nil return makes
// it directly usable with promhttp.WithExemplarFromContext.
func ExemplarFromContext(ctx context.Context) prometheus.Labels {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() || !sc.IsSampled() {
		return nil
	}
	return prometheus.Labels{
		"trace_id": sc.TraceID().String(),
		"span_id":  sc.SpanID().String(),
	}
}

// AddWithExemplar adds v to counter. When ctx carries a sampled trace span,
// the addition carries a trace exemplar so dashboards can link the counter
// to the trace. Any panic from exemplar validation is swallowed: a lost
// exemplar must never break the recording path, and client_golang applies
// the addition before validating the exemplar, so the count is recorded
// either way.
func AddWithExemplar(ctx context.Context, counter prometheus.Counter, v float64) {
	exemplar := ExemplarFromContext(ctx)
	adder, ok := counter.(prometheus.ExemplarAdder)
	if exemplar == nil || !ok {
		counter.Add(v)
		return
	}
	defer func() {
		_ = recover()
	}()
	adder.AddWithExemplar(v, exemplar)
}

// ObserveWithExemplar records v on obs. When ctx carries a sampled trace
// span, the observation carries a trace exemplar on the histogram bucket it
// falls into. The same panic-swallowing rationale as AddWithExemplar
// applies.
func ObserveWithExemplar(ctx context.Context, obs prometheus.Observer, v float64) {
	exemplar := ExemplarFromContext(ctx)
	observer, ok := obs.(prometheus.ExemplarObserver)
	if exemplar == nil || !ok {
		obs.Observe(v)
		return
	}
	defer func() {
		_ = recover()
	}()
	observer.ObserveWithExemplar(v, exemplar)
}
