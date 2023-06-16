// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package events

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ory/x/otelx/semconv"
	otelattr "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	OAuth2ClientCreated          semconv.Event = "OAuth2ClientCreated"
	OAuth2ClientUpdated          semconv.Event = "OAuth2ClientUpdated"
	OAuth2ClientDeleted          semconv.Event = "OAuth2ClientDeleted"
	OAuth2ConsentRequestAccepted semconv.Event = "OAuth2ConsentRequestAccepted"
	OAuth2ConsentRequestRejected semconv.Event = "OAuth2ConsentRequestRejected"
	OAuth2LoginRequestAccepted   semconv.Event = "OAuth2LoginRequestAccepted"
	OAuth2LoginRequestRejected   semconv.Event = "OAuth2LoginRequestRejected"
	OAuth2AccessTokenIssued      semconv.Event = "OAuth2AccessTokenIssued"
	OAuth2RefreshTokenIssued     semconv.Event = "OAuth2RefreshTokenIssued"
	OIDCIDTokenIssued            semconv.Event = "OIDCIDTokenIssued"
)

const (
	attributeKeyOAuth2ClientID semconv.AttributeKey = "OAuth2ClientID"
	attributeKeyOAuth2Subject  semconv.AttributeKey = "OAuth2Subject"
)

func attrOAuth2ClientID(val uuid.UUID) otelattr.KeyValue {
	return otelattr.String(attributeKeyOAuth2ClientID.String(), val.String())
}

func attrOAuth2Subject(sub string) otelattr.KeyValue {
	return otelattr.String(attributeKeyOAuth2Subject.String(), sub)
}

// NewOAuth2ClientCreated returns a new OAuth2ClientCreated event.
func NewOAuth2ClientCreated(ctx context.Context, clientID uuid.UUID) (string, trace.EventOption) {
	return OAuth2ClientCreated.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
			)...,
		)
}

// NewOAuth2ClientUpdated returns a new OAuth2ClientUpdated event.
func NewOAuth2ClientUpdated(ctx context.Context, clientID uuid.UUID) (string, trace.EventOption) {
	return OAuth2ClientUpdated.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
			)...,
		)
}

// NewOAuth2ClientDeleted returns a new OAuth2ClientDeleted event.
func NewOAuth2ClientDeleted(ctx context.Context, clientID uuid.UUID) (string, trace.EventOption) {
	return OAuth2ClientDeleted.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
			)...,
		)
}

// NewOAuth2ConsentRequestAccepted returns a new OAuth2ConsentRequestAccepted event.
func NewOAuth2ConsentRequestAccepted(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OAuth2ConsentRequestAccepted.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}

// NewOAuth2ConsentRequestRejected returns a new OAuth2ConsentRequestRejected event.
func NewOAuth2ConsentRequestRejected(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OAuth2ConsentRequestRejected.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}

// NewOAuth2LoginRequestAccepted returns a new OAuth2LoginRequestAccepted event.
func NewOAuth2LoginRequestAccepted(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OAuth2LoginRequestAccepted.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}

// NewOAuth2LoginRequestRejected returns a new OAuth2LoginRequestRejected event.
func NewOAuth2LoginRequestRejected(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OAuth2LoginRequestRejected.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}

// NewOAuth2AccessTokenIssued returns a new OAuth2AccessTokenIssued event.
func NewOAuth2AccessTokenIssued(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OAuth2AccessTokenIssued.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}

// NewOAuth2RefreshTokenIssued returns a new OAuth2RefreshTokenIssued event.
func NewOAuth2RefreshTokenIssued(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OAuth2RefreshTokenIssued.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}

// NewOIDCIDTokenIssued returns a new OIDCIDTokenIssued event.
func NewOIDCIDTokenIssued(ctx context.Context, clientID uuid.UUID, subject string) (string, trace.EventOption) {
	return OIDCIDTokenIssued.String(),
		trace.WithAttributes(
			append(
				semconv.AttributesFromContext(ctx),
				attrOAuth2ClientID(clientID),
				attrOAuth2Subject(subject),
			)...,
		)
}
