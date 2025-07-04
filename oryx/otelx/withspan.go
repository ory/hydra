// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package otelx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

// WithSpan wraps execution of f in a span identified by name.
//
// If f returns an error or panics, the span status will be set to the error
// state. The error (or panic) will be propagated unmodified.
//
// f will be wrapped in a child span by default. To make a new root span
// instead, pass the trace.WithNewRoot() option.
func WithSpan(ctx context.Context, name string, f func(context.Context) error, opts ...trace.SpanStartOption) (err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, name, opts...)
	defer func() {
		defer span.End()
		if r := recover(); r != nil {
			setErrorStatusPanic(span, r)
			panic(r)
		} else if err != nil {
			span.SetStatus(codes.Error, err.Error())
			setErrorTags(span, err)
		}
	}()
	return f(ctx)
}

// End finishes span, and automatically sets the error state if *err is not nil
// or during panicking.
//
// Usage:
//
//	func Divide(ctx context.Context, numerator, denominator int) (ratio int, err error) {
//		ctx, span := tracer.Start(ctx, "Divide")
//		defer otelx.End(span, &err)
//		if denominator == 0 {
//			return 0, errors.New("cannot divide by zero")
//		}
//		return numerator / denominator, nil
//	}
//
// During a panic, we don't fully conform to OpenTelemetry's semantic
// conventions because that would require us to emit a span event to attach the
// stacktrace and error type, and we don't want to do that. Instead, we set the
// tags on the span directly.
// https://opentelemetry.io/docs/specs/semconv/exceptions/exceptions-spans/
//
// For improved compatibility with Datadog, we also set some additional tags as
// documented here:
// https://docs.datadoghq.com/standard-attributes/?product=apm&search=error
func End(span trace.Span, err *error) {
	defer span.End()
	if r := recover(); r != nil {
		setErrorStatusPanic(span, r)
		panic(r)
	}
	if err == nil || *err == nil {
		return
	}
	span.SetStatus(codes.Error, (*err).Error())
	setErrorTags(span, *err)
}

func setErrorStatusPanic(span trace.Span, recovered any) {
	span.SetAttributes(
		// OpenTelemetry says to add these attributes to an event, not the span
		// itself. We don't want to do that, so we're adding them to the span
		// directly.
		semconv.ExceptionEscaped(true),
		// OpenTelemetry describes "exception.stacktrace"  We don't love that,
		// though, so we're using "error.stack" instead, like DataDog).
		attribute.String("error.stack", stacktrace()),
	)
	if t := reflect.TypeOf(recovered); t != nil {
		span.SetAttributes(semconv.ExceptionType(t.String()))
	}
	switch e := recovered.(type) {
	case error:
		span.SetStatus(codes.Error, "panic: "+e.Error())
		setErrorTags(span, e)
	case string, fmt.Stringer:
		span.SetStatus(codes.Error, fmt.Sprintf("panic: %v", e))
	default:
		span.SetStatus(codes.Error, "panic")
	case nil:
		// nothing
	}
}

func setErrorTags(span trace.Span, err error) {
	span.SetAttributes(
		attribute.String("error", err.Error()),
		attribute.String("error.message", err.Error()),                        // DataDog compat
		attribute.String("error.type", fmt.Sprintf("%T", errors.Unwrap(err))), // the innermost error type is the most useful here
	)
	if e := interface{ StackTrace() pkgerrors.StackTrace }(nil); errors.As(err, &e) {
		span.SetAttributes(attribute.String("error.stack", fmt.Sprintf("%+v", e.StackTrace())))
	}
	if e := interface{ Reason() string }(nil); errors.As(err, &e) {
		span.SetAttributes(attribute.String("error.reason", e.Reason()))
	}
	if e := interface{ Debug() string }(nil); errors.As(err, &e) {
		span.SetAttributes(attribute.String("error.debug", e.Debug()))
	}
	if e := interface{ ID() string }(nil); errors.As(err, &e) {
		span.SetAttributes(attribute.String("error.id", e.ID()))
	}
	if e := interface{ Details() map[string]interface{} }(nil); errors.As(err, &e) {
		for k, v := range e.Details() {
			span.SetAttributes(attribute.String("error.details."+k, fmt.Sprintf("%v", v)))
		}
	}
}

func stacktrace() string {
	pc := make([]uintptr, 5)
	n := runtime.Callers(4, pc)
	if n == 0 {
		return ""
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return builder.String()
}
