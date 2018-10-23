package instrumentedsql

import "context"

// Tracer is the interface needed to be implemented by any tracing implementation we use
type Tracer interface {
	GetSpan(ctx context.Context) Span
}

// Span is part of the interface needed to be implemented by any tracing implementation we use
type Span interface {
	NewChild(string) Span
	SetLabel(k, v string)
	SetError(err error)
	Finish()
}

type nullTracer struct{}
type nullSpan struct{}

func (nullTracer) GetSpan(ctx context.Context) Span {
	return nullSpan{}
}

func (nullSpan) NewChild(string) Span {
	return nullSpan{}
}

func (nullSpan) SetLabel(k, v string) {}

func (nullSpan) Finish() {}

func (nullSpan) SetError(err error) {}
