// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package semconv contains OpenTelemetry semantic convention constants.
package semconv

import (
	"github.com/gofrs/uuid"
	otelattr "go.opentelemetry.io/otel/attribute"

	"github.com/ory/x/httpx"
)

type Event string

func (e Event) String() string {
	return string(e)
}

type AttributeKey string

func (a AttributeKey) String() string {
	return string(a)
}

const (
	AttributeKeyIdentityID         AttributeKey = "IdentityID"
	AttributeKeyNID                AttributeKey = "ProjectID"
	AttributeKeyClientIP           AttributeKey = "ClientIP"
	AttributeKeyGeoLocationCity    AttributeKey = "GeoLocationCity"
	AttributeKeyGeoLocationRegion  AttributeKey = "GeoLocationRegion"
	AttributeKeyGeoLocationCountry AttributeKey = "GeoLocationCountry"
	AttributeKeyWorkspace          AttributeKey = "WorkspaceID"
	AttributeKeySubscriptionID     AttributeKey = "SubscriptionID"
	AttributeKeyProjectEnvironment AttributeKey = "ProjectEnvironment"
	AttributeKeyWorkspaceAPIKeyID  AttributeKey = "WorkspaceAPIKeyID"
	AttributeKeyProjectAPIKeyID    AttributeKey = "ProjectAPIKeyID"
)

func AttrIdentityID[V string | uuid.UUID](val V) otelattr.KeyValue {
	return otelattr.String(AttributeKeyIdentityID.String(), uuidOrString(val))
}

func AttrNID(val uuid.UUID) otelattr.KeyValue {
	return otelattr.String(AttributeKeyNID.String(), val.String())
}

func AttrWorkspace(val uuid.UUID) otelattr.KeyValue {
	return otelattr.String(AttributeKeyWorkspace.String(), val.String())
}

func AttrSubscription(val uuid.UUID) otelattr.KeyValue {
	return otelattr.String(AttributeKeySubscriptionID.String(), val.String())
}

func AttrProjectEnvironment(val string) otelattr.KeyValue {
	return otelattr.String(AttributeKeyProjectEnvironment.String(), val)
}

func AttrClientIP(val string) otelattr.KeyValue {
	return otelattr.String(AttributeKeyClientIP.String(), val)
}

func AttrGeoLocation(val httpx.GeoLocation) []otelattr.KeyValue {
	geoLocationAttributes := make([]otelattr.KeyValue, 0, 3)

	if val.City != "" {
		geoLocationAttributes = append(geoLocationAttributes, otelattr.String(AttributeKeyGeoLocationCity.String(), val.City))
	}
	if val.Country != "" {
		geoLocationAttributes = append(geoLocationAttributes, otelattr.String(AttributeKeyGeoLocationCountry.String(), val.Country))
	}
	if val.Region != "" {
		geoLocationAttributes = append(geoLocationAttributes, otelattr.String(AttributeKeyGeoLocationRegion.String(), val.Region))
	}

	return geoLocationAttributes
}

func AttrWorkspaceAPIKeyID[V string | uuid.UUID](val V) otelattr.KeyValue {
	return otelattr.String(AttributeKeyWorkspaceAPIKeyID.String(), uuidOrString(val))
}

func AttrProjectAPIKeyID[V string | uuid.UUID](val V) otelattr.KeyValue {
	return otelattr.String(AttributeKeyProjectAPIKeyID.String(), uuidOrString(val))
}

func uuidOrString[V string | uuid.UUID](val V) string {
	switch val := any(val).(type) {
	case string:
		return val
	case uuid.UUID:
		return val.String()
	}
	panic("unreachable")
}
