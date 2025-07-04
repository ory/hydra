// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package semconv

import (
	"context"

	otelattr "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewDeprecatedFeatureUsedEvent creates a new event indicating that a deprecated feature was used.
// It returns the event name and a trace.EventOption that can be used to
// add the event to a span.
//
//	span.AddEvent(NewDeprecatedFeatureUsedEvent(ctx, "deprecated-feature-id", otelattr.String("key", "value")))
func NewDeprecatedFeatureUsedEvent(ctx context.Context, deprecatedCodeFeatureID string, attrs ...otelattr.KeyValue) (string, trace.EventOption) {
	return DeprecatedFeatureUsed.String(),
		trace.WithAttributes(
			append(
				append(
					attrs,
					AttributesFromContext(ctx)...,
				),
				AttrDeprecatedFeatureID(deprecatedCodeFeatureID),
			)...,
		)
}

const (
	AttributeKeyDeprecatedCodePathIDAttributeKey AttributeKey = "DeprecatedFeatureID"
	DeprecatedFeatureUsed                        Event        = "DeprecatedFeatureUsed"
)

func AttrDeprecatedFeatureID(id string) otelattr.KeyValue {
	return otelattr.String(AttributeKeyDeprecatedCodePathIDAttributeKey.String(), id)
}
