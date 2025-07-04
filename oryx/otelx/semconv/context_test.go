// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package semconv

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"

	"github.com/ory/x/httpx"
)

func TestAttributesFromContext(t *testing.T) {
	ctx := context.Background()
	assert.Len(t, AttributesFromContext(ctx), 0)

	nid, wsID := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	ctx = ContextWithAttributes(ctx, AttrNID(nid), AttrWorkspace(wsID))
	assert.Len(t, AttributesFromContext(ctx), 2)

	uid1, uid2 := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	location := httpx.GeoLocation{
		City:    "Berlin",
		Country: "Germany",
		Region:  "BE",
	}
	ctx = ContextWithAttributes(ctx, append(AttrGeoLocation(location), AttrIdentityID(uid1), AttrClientIP("127.0.0.1"), AttrIdentityID(uid2))...)
	attrs := AttributesFromContext(ctx)
	assert.Len(t, attrs, 7, "should deduplicate")
	assert.Equal(t, []attribute.KeyValue{
		attribute.String(AttributeKeyNID.String(), nid.String()),
		attribute.String(AttributeKeyWorkspace.String(), wsID.String()),
		attribute.String(AttributeKeyGeoLocationCity.String(), "Berlin"),
		attribute.String(AttributeKeyGeoLocationCountry.String(), "Germany"),
		attribute.String(AttributeKeyGeoLocationRegion.String(), "BE"),
		attribute.String(AttributeKeyClientIP.String(), "127.0.0.1"),
		attribute.String(AttributeKeyIdentityID.String(), uid2.String()),
	}, attrs, "last duplicate attribute wins")
}
