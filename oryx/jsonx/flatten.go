// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Flatten flattens a JSON object using dot notation.
func Flatten(raw json.RawMessage) map[string]interface{} {
	parsed := gjson.ParseBytes(raw)
	if !parsed.IsObject() {
		return nil
	}

	flattened := make(map[string]interface{})
	flatten(parsed, nil, flattened)
	return flattened
}

func flatten(parsed gjson.Result, parents []string, flattened map[string]interface{}) {
	if parsed.IsObject() {
		parsed.ForEach(func(k, v gjson.Result) bool {
			flatten(v, append(parents, strings.ReplaceAll(k.String(), ".", "\\.")), flattened)
			return true
		})
	} else if parsed.IsArray() {
		for kk, vv := range parsed.Array() {
			flatten(vv, append(parents, strconv.Itoa(kk)), flattened)
		}
	} else {
		flattened[strings.Join(parents, ".")] = parsed.Value()
	}
}
