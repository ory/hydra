package opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/luna-duclos/instrumentedsql"
	"github.com/opentracing/opentracing-go"
)

type tracer struct {
	traceOrphans bool
}

type span struct {
	tracer
	parent opentracing.Span
}

// NewTracer returns a tracer that will fetch spans using opentracing's SpanFromContext function
// if traceOrphans is set to true, then spans with no parent will be traced anyway, if false, they will not be.
func NewTracer(traceOrphans bool) instrumentedsql.Tracer { return tracer{traceOrphans: traceOrphans} }

// GetSpan returns a span
func (t tracer) GetSpan(ctx context.Context) instrumentedsql.Span {
	if ctx == nil {
		return span{parent: nil, tracer: t}
	}

	return span{parent: opentracing.SpanFromContext(ctx), tracer: t}
}

func (s span) NewChild(name string) instrumentedsql.Span {
	if s.parent == nil {
		if s.traceOrphans {
			return span{parent: opentracing.StartSpan(name), tracer: s.tracer}
		}

		return s
	}

	return span{parent: opentracing.StartSpan(name, opentracing.ChildOf(s.parent.Context())), tracer: s.tracer}
}

func (s span) SetLabel(k, v string) {
	if s.parent == nil {
		return
	}
	s.parent.SetTag(k, v)
}

func (s span) SetError(err error) {
	if err == nil {
		return
	}

	if s.parent == nil {
		return
	}

	ext.Error.Set(s.parent, true)
}

func (s span) Finish() {
	if s.parent == nil {
		return
	}
	s.parent.Finish()
}
