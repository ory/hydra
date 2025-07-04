// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package semconv

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/ory/x/httpx"
)

type contextKey int

const contextKeyAttributes contextKey = iota

func ContextWithAttributes(ctx context.Context, attrs ...attribute.KeyValue) context.Context {
	existing, _ := ctx.Value(contextKeyAttributes).([]attribute.KeyValue)
	return context.WithValue(ctx, contextKeyAttributes, append(existing, attrs...))
}

func AttributesFromContext(ctx context.Context) []attribute.KeyValue {
	fromCtx, _ := ctx.Value(contextKeyAttributes).([]attribute.KeyValue)
	uniq := make(map[attribute.Key]struct{})
	attrs := make([]attribute.KeyValue, 0)
	for i := len(fromCtx) - 1; i >= 0; i-- {
		if _, ok := uniq[fromCtx[i].Key]; !ok {
			uniq[fromCtx[i].Key] = struct{}{}
			attrs = append(attrs, fromCtx[i])
		}
	}
	reverse(attrs)
	return attrs
}

func Middleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := ContextWithAttributes(r.Context(),
		append(
			AttrGeoLocation(*httpx.ClientGeoLocation(r)),
			AttrClientIP(httpx.ClientIP(r)),
		)...,
	)

	next(rw, r.WithContext(ctx))
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
