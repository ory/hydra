// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package mapx

import (
	"encoding/json"
	"errors"
	"math"
	"time"
)

// ErrKeyDoesNotExist is returned when the key does not exist in the map.
var ErrKeyDoesNotExist = errors.New("key is not present in map")

// ErrKeyCanNotBeTypeAsserted is returned when the key can not be type asserted.
var ErrKeyCanNotBeTypeAsserted = errors.New("key could not be type asserted")

// GetString returns a string for a given key in values.
func GetString[K comparable](values map[K]any, key K) (string, error) {
	if v, ok := values[key]; !ok {
		return "", ErrKeyDoesNotExist
	} else if sv, ok := v.(string); !ok {
		return "", ErrKeyCanNotBeTypeAsserted
	} else {
		return sv, nil
	}
}

// GetStringSlice returns a string slice for a given key in values.
func GetStringSlice[K comparable](values map[K]any, key K) ([]string, error) {
	v, ok := values[key]
	if !ok {
		return nil, ErrKeyDoesNotExist
	}
	switch v := v.(type) {
	case []string:
		return v, nil
	case []any:
		vs := make([]string, len(v))
		for k, v := range v {
			var ok bool
			vs[k], ok = v.(string)
			if !ok {
				return nil, ErrKeyCanNotBeTypeAsserted
			}
		}
		return vs, nil
	}
	return nil, ErrKeyCanNotBeTypeAsserted
}

// GetTime returns a string slice for a given key in values.
func GetTime[K comparable](values map[K]any, key K) (time.Time, error) {
	v, ok := values[key]
	if !ok {
		return time.Time{}, ErrKeyDoesNotExist
	}

	switch v := v.(type) {
	case time.Time:
		return v, nil
	case int64:
		return time.Unix(v, 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	case float64:
		if v < math.MinInt64 || v > math.MaxInt64 {
			return time.Time{}, errors.New("value is out of range")
		}
		return time.Unix(int64(v), 0), nil
	case float32:
		if v < math.MinInt64 || v > math.MaxInt64 {
			return time.Time{}, errors.New("value is out of range")
		}
		return time.Unix(int64(v), 0), nil
	}

	return time.Time{}, ErrKeyCanNotBeTypeAsserted
}

// GetInt64 returns an int64 for a given key in values.
func GetInt64[K comparable](values map[K]any, key K) (int64, error) {
	v, ok := values[key]
	if !ok {
		return 0, ErrKeyDoesNotExist
	}
	switch v := v.(type) {
	case json.Number:
		return v.Int64()
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint:
		vv := uint64(v)
		if vv > math.MaxInt64 {
			return 0, errors.New("value is out of range")
		}
		return int64(vv), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, errors.New("value is out of range")
		}
		return int64(v), nil
	}
	return 0, ErrKeyCanNotBeTypeAsserted
}

// GetInt32 returns an int32 for a given key in values.
func GetInt32[K comparable](values map[K]any, key K) (int32, error) {
	v, err := GetInt64(values, key)
	if err != nil {
		return 0, err
	}
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, errors.New("value is out of range")
	}
	return int32(v), nil
}

// GetInt returns an int for a given key in values.
func GetInt[K comparable](values map[K]any, key K) (int, error) {
	v, err := GetInt64(values, key)
	if err != nil {
		return 0, err
	}
	if v > math.MaxInt || v < math.MinInt {
		return 0, errors.New("value is out of range")
	}
	return int(v), nil
}

// GetFloat64Default returns a float64 or the default value for a given key in values.
func GetFloat64Default[K comparable](values map[K]any, key K, defaultValue float64) float64 {
	f, err := GetFloat64(values, key)
	if err != nil {
		return defaultValue
	}
	return f
}

// GetFloat64 returns a float64 for a given key in values.
func GetFloat64[K comparable](values map[K]any, key K) (float64, error) {
	v, ok := values[key]
	if !ok {
		return 0, ErrKeyDoesNotExist
	}
	switch v := v.(type) {
	case json.Number:
		return v.Float64()
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, ErrKeyCanNotBeTypeAsserted
}

// GetStringDefault returns a string or the default value for a given key in values.
func GetStringDefault[K comparable](values map[K]any, key K, defaultValue string) string {
	if s, err := GetString(values, key); err == nil {
		return s
	}
	return defaultValue
}

// GetStringSliceDefault returns a string slice or the default value for a given key in values.
func GetStringSliceDefault[K comparable](values map[K]any, key K, defaultValue []string) []string {
	if s, err := GetStringSlice(values, key); err == nil {
		return s
	}
	return defaultValue
}
