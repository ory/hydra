// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package otelx

import (
	"database/sql"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

const nullString = "<null>"

func StringAttrs(attrs map[string]string) []attribute.KeyValue {
	s := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		s = append(s, attribute.String(k, v))
	}
	return s
}

func AutoInt[I int | int32 | int64](k string, v I) attribute.KeyValue {
	// Internally, the OpenTelemetry SDK uses int64 for all integer values anyway.
	return attribute.Int64(k, int64(v))
}

func Nullable[V any, VN *V | sql.Null[V], A func(string, V) attribute.KeyValue](a A, k string, v VN) attribute.KeyValue {
	switch v := any(v).(type) {
	case *V:
		if v == nil {
			return attribute.String(k, nullString)
		}
		return a(k, *v)
	case sql.Null[V]:
		if !v.Valid {
			return attribute.String(k, nullString)
		}
		return a(k, v.V)
	}
	// This should never happen, as the type switch above is exhaustive to the generic type VN.
	return attribute.String(k, fmt.Sprintf("<got unsupported type %T>", v))
}

func NullString[V *string | sql.Null[string]](k string, v V) attribute.KeyValue {
	return Nullable(attribute.String, k, v)
}

func NullStringer(k string, v fmt.Stringer) attribute.KeyValue {
	if v == nil {
		return attribute.String(k, nullString)
	}
	return attribute.String(k, v.String())
}
