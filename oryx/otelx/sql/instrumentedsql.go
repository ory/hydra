// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"database/sql/driver"

	"github.com/luna-duclos/instrumentedsql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracingComponent = "github.com/ory/x/otelx/sql"

type (
	tracer struct{}
	span   struct {
		ctx    context.Context
		parent trace.Span
	}
)

var (
	_ instrumentedsql.Tracer = tracer{}
	_ instrumentedsql.Span   = span{}
)

func NewTracer() instrumentedsql.Tracer { return tracer{} }

// GetSpan returns a span
func (tracer) GetSpan(ctx context.Context) instrumentedsql.Span {
	return span{ctx, trace.SpanFromContext(ctx)}
}

func (s span) NewChild(name string) instrumentedsql.Span {
	ctx, child := s.parent.TracerProvider().Tracer(tracingComponent).Start(s.ctx, name, trace.WithSpanKind(trace.SpanKindClient))
	return span{ctx, child}
}

func (s span) SetLabel(k, v string) {
	s.parent.SetAttributes(attribute.String(k, v))
}

func (s span) SetError(err error) {
	if err == nil || err == driver.ErrSkip {
		return
	}
	s.parent.SetStatus(codes.Error, err.Error())
}

func (s span) Finish() {
	s.parent.End()
}
