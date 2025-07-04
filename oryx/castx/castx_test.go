// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package castx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFloatSliceE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect []float64
		iserr  bool
	}{
		{[]int{1, 3}, []float64{1, 3}, false},
		{[]interface{}{1.2, 3.2}, []float64{1.2, 3.2}, false},
		{[]string{"2", "3"}, []float64{2, 3}, false},
		{[]string{"2.2", "3.2"}, []float64{2.2, 3.2}, false},
		{[2]string{"2", "3"}, []float64{2, 3}, false},
		{[2]string{"2.2", "3.2"}, []float64{2.2, 3.2}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToFloatSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errmsg)
			continue
		}

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToFloatSlice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	}
}

func TestToStringSlice(t *testing.T) {
	assert.Equal(t, []string{"foo", "bar"}, ToStringSlice("foo,bar"))
	assert.NotEqual(t, []string{"foo bar baz"}, ToStringSlice("foo bar baz,"))
	assert.Equal(t, []string{"foo bar baz", ""}, ToStringSlice("foo bar baz,"))
	assert.NotEqual(t, []string{"foo", "bar", "baz"}, ToStringSlice("foo bar baz"))
	assert.Equal(t, []string{"foo bar baz"}, ToStringSlice("foo bar baz"))
	assert.Equal(t, []string{"foo", "bar", "baz,", " asdf"}, ToStringSlice("foo,bar,\"baz,\", asdf"))
	assert.Equal(t, []string{"'foo'", "x\"bar", "baz"}, ToStringSlice("'foo',\"x\"\"bar\",baz"))
}
