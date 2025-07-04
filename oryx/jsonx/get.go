// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"reflect"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func jsonKey(f reflect.StructField) *string {
	if jsonTag := f.Tag.Get("json"); jsonTag != "" {
		if jsonTag == "-" {
			return nil
		}
		return &strings.Split(jsonTag, ",")[0]
	} else if f.Anonymous {
		return nil
	} else if f.IsExported() {
		return &f.Name
	}
	return nil
}

// AllValidJSONKeys returns all JSON keys from the struct or *struct type.
// It does not return keys from nested slices, but embedded/nested structs.
func AllValidJSONKeys(s interface{}) (keys []string) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := range t.NumField() {
		f := t.Field(i)
		jKey := jsonKey(f)
		if k := f.Type.Kind(); k == reflect.Struct || k == reflect.Ptr {
			subKeys := AllValidJSONKeys(v.Field(i).Interface())
			for _, subKey := range subKeys {
				if jKey != nil {
					keys = append(keys, *jKey+"."+subKey)
				} else {
					keys = append(keys, subKey)
				}
			}
		} else if jKey != nil {
			keys = append(keys, *jKey)
		}
	}
	return keys
}

// ParseEnsureKeys returns a result that has the GetRequireValidKey function.
func ParseEnsureKeys(original interface{}, raw []byte) *Result {
	return &Result{
		keys:   AllValidJSONKeys(original),
		result: gjson.ParseBytes(raw),
	}
}

type Result struct {
	result gjson.Result
	keys   []string
}

// GetRequireValidKey ensures that the key is valid before returning the result.
func (r *Result) GetRequireValidKey(t require.TestingT, key string) gjson.Result {
	require.Contains(t, r.keys, key)
	return r.result.Get(key)
}

func GetRequireValidKey(t require.TestingT, original interface{}, raw []byte, key string) gjson.Result {
	return ParseEnsureKeys(original, raw).GetRequireValidKey(t, key)
}
