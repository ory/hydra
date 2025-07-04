// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package semconv

import (
	"context"

	otelattr "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewWarning creates a new warning event with the given ID and attributes.
// It returns the event name and a trace.EventOption that can be used to
// add the event to a span.
//
//	span.AddEvent(NewWarning(ctx, "warning-id", otelattr.String("key", "value")))
func NewWarning(ctx context.Context, id string, attrs ...otelattr.KeyValue) (string, trace.EventOption) {
	return Warning.String(),
		trace.WithAttributes(
			append(
				append(
					attrs,
					AttributesFromContext(ctx)...,
				),
				otelattr.String(AttributeWarningID.String(), id),
			)...,
		)
}

const (
	Warning            Event        = "Warning"
	AttributeWarningID AttributeKey = "WarningID"
)

func AttrWarningID(id string) otelattr.KeyValue {
	return otelattr.String(AttributeWarningID.String(), id)
}
