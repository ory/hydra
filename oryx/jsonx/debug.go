// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"encoding/json"
	"fmt"
	"slices"
)

// Anonymize takes a JSON byte array and anonymizes its content by
// recursively replacing all values with a string indicating their type.
//
// It recurses into nested objects and arrays, but ignores the "schemas" and "id".
func Anonymize(data []byte, except ...string) []byte {
	obj := make(map[string]any)
	if err := json.Unmarshal(data, &obj); err != nil {
		return []byte(fmt.Sprintf(`{"error": "invalid JSON", "message": %q}`, err.Error()))
	}

	anonymize(obj, except...)

	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return []byte(fmt.Sprintf(`{"error": "could not marshal JSON shape", "message": %q}`, err.Error()))
	}

	return out
}

func anonymize(obj map[string]any, except ...string) {
	for k, v := range obj {
		if slices.Contains(except, k) {
			continue
		}

		switch v := v.(type) {
		case []any:
			for elIdx, el := range v {
				switch el := el.(type) {
				case map[string]any:
					anonymize(el)
					v[elIdx] = el
				default:
					v[elIdx] = jsonType(el)
				}
			}

		case map[string]any:
			anonymize(v)
			obj[k] = v
		default:
			obj[k] = jsonType(v)
		}
	}
}

func jsonType(v any) string {
	switch v := v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return fmt.Sprintf("%T", v)
	}
}
