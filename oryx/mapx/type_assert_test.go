// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package mapx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetString(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": 1234}
	v, err := GetString(m, "foo")
	require.NoError(t, err)
	assert.EqualValues(t, "bar", v)
	_, err = GetString(m, "bar")
	require.Error(t, err)
	_, err = GetString(m, "baz")
	require.Error(t, err)
}

func TestGetStringSlice(t *testing.T) {
	m := map[interface{}]interface{}{"foo": []string{"foo", "bar"}, "baz": "bar"}
	v, err := GetStringSlice(m, "foo")
	require.NoError(t, err)
	assert.EqualValues(t, []string{"foo", "bar"}, v)
	_, err = GetStringSlice(m, "bar")
	require.Error(t, err)
	_, err = GetStringSlice(m, "baz")
	require.Error(t, err)
}

func TestGetStringSliceDefault(t *testing.T) {
	m := map[interface{}]interface{}{"foo": []string{"foo", "bar"}, "baz": "bar"}
	assert.EqualValues(t, []string{"foo", "bar"}, GetStringSliceDefault(m, "foo", []string{"default"}))
	assert.EqualValues(t, []string{"default"}, GetStringSliceDefault(m, "baz", []string{"default"}))
	assert.EqualValues(t, []string{"default"}, GetStringSliceDefault(m, "bar", []string{"default"}))
}

func TestGetStringDefault(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": 1234}
	assert.EqualValues(t, "bar", GetStringDefault(m, "foo", "default"))
	assert.EqualValues(t, "default", GetStringDefault(m, "baz", "default"))
	assert.EqualValues(t, "default", GetStringDefault(m, "bar", "default"))
}

func TestGetFloat32(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": float32(1234)}
	v, err := GetFloat32(m, "baz")
	require.NoError(t, err)
	assert.EqualValues(t, float32(1234), v)
	_, err = GetFloat32(m, "foo")
	require.Error(t, err)
	_, err = GetFloat32(m, "bar")
	require.Error(t, err)
}

func TestGetFloat64(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": float64(1234)}
	v, err := GetFloat64(m, "baz")
	require.NoError(t, err)
	assert.EqualValues(t, float64(1234), v)
	_, err = GetFloat64(m, "foo")
	require.Error(t, err)
	_, err = GetFloat64(m, "bar")
	require.Error(t, err)
}

func TestGetGetFloat64Default(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": float64(1234)}
	v := GetFloat64Default(m, "baz", 0)
	assert.EqualValues(t, float64(1234), v)
	v = GetFloat64Default(m, "foo", float64(1))
	assert.EqualValues(t, float64(1), v)
	v = GetFloat64Default(m, "bar", float64(2))
	assert.EqualValues(t, float64(2), v)
}

func TestGetGetFloat32Default(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": float32(1234)}
	v := GetFloat32Default(m, "baz", 0)
	assert.EqualValues(t, float32(1234), v)
	v = GetFloat32Default(m, "foo", float32(1))
	assert.EqualValues(t, float32(1), v)
	v = GetFloat32Default(m, "bar", float32(2))
	assert.EqualValues(t, float32(2), v)
}

func TestGetGetInt32Default(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": int32(1234)}
	v := GetInt32Default(m, "baz", 0)
	assert.EqualValues(t, int32(1234), v)
	v = GetInt32Default(m, "foo", int32(1))
	assert.EqualValues(t, int32(1), v)
	v = GetInt32Default(m, "bar", int32(2))
	assert.EqualValues(t, int32(2), v)
}

func TestGetGetInt64Default(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": int64(1234)}
	v := GetInt64Default(m, "baz", 0)
	assert.EqualValues(t, int64(1234), v)
	v = GetInt64Default(m, "foo", int64(1))
	assert.EqualValues(t, int64(1), v)
	v = GetInt64Default(m, "bar", int64(2))
	assert.EqualValues(t, int64(2), v)
}

func TestGetGetIntDefault(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": int(1234)}
	v := GetIntDefault(m, "baz", 0)
	assert.EqualValues(t, int(1234), v)
	v = GetIntDefault(m, "foo", int(1))
	assert.EqualValues(t, int(1), v)
	v = GetIntDefault(m, "bar", int(2))
	assert.EqualValues(t, int(2), v)
}

func TestGetInt64(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": int64(1234)}
	v, err := GetInt64(m, "baz")
	require.NoError(t, err)
	assert.EqualValues(t, int64(1234), v)
	_, err = GetInt64(m, "foo")
	require.Error(t, err)
	_, err = GetInt64(m, "bar")
	require.Error(t, err)
}

func TestGetInt32(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": int32(1234), "baz2": int(1234)}
	v, err := GetInt32(m, "baz")
	require.NoError(t, err)
	assert.EqualValues(t, int32(1234), v)
	v, err = GetInt32(m, "baz2")
	require.NoError(t, err)
	assert.EqualValues(t, int32(1234), v)
	_, err = GetInt32(m, "foo")
	require.Error(t, err)
	_, err = GetInt32(m, "bar")
	require.Error(t, err)
}

func TestKeyStringToInterface(t *testing.T) {
	assert.EqualValues(t, map[interface{}]interface{}{"foo": "bar", "baz": 1234, "baz2": int32(1234)}, KeyStringToInterface(map[string]interface{}{"foo": "bar", "baz": 1234, "baz2": int32(1234)}))
}

func TestGetInt(t *testing.T) {
	m := map[interface{}]interface{}{"foo": "bar", "baz": 1234, "baz2": int32(1234)}
	v, err := GetInt32(m, "baz")
	require.NoError(t, err)
	assert.EqualValues(t, int32(1234), v)
	_, err = GetInt32(m, "foo")
	require.Error(t, err)
	_, err = GetInt32(m, "bar")
	require.Error(t, err)
}

func TestToJSONMap(t *testing.T) {
	assert.EqualValues(t, map[string]interface{}{"baz": []interface{}{map[string]interface{}{"bar": "bar"}}, "foo": "bar"}, ToJSONMap(map[string]interface{}{
		"foo": "bar",
		"baz": []interface{}{
			map[interface{}]interface{}{
				"bar": "bar",
			},
		},
	}))

}
