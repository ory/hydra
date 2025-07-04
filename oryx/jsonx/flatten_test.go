// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatten(t *testing.T) {
	f, err := os.ReadFile("./stub/random.json")
	require.NoError(t, err)

	for k, tc := range []struct {
		raw      []byte
		expected map[string]interface{}
	}{
		{
			raw:      f,
			expected: map[string]interface{}{"fall": "to", "floating.0": -1.273085434e+09, "floating.1": 9.53442581e+08, "floating.2.gray.buy": true, "floating.2.gray.hold.0.0": 1.81518765e+08, "floating.2.gray.hold.0.1.0.flies": -1.571371799e+09, "floating.2.gray.hold.0.1.0.leather": "across", "floating.2.gray.hold.0.1.0.over": 5.12666854e+08, "floating.2.gray.hold.0.1.0.shaking": true, "floating.2.gray.hold.0.1.0.steam.ago": true, "floating.2.gray.hold.0.1.0.steam.appropriate": 1.249911539e+09, "floating.2.gray.hold.0.1.0.steam.box": false, "floating.2.gray.hold.0.1.0.steam.cry": 1.463961818e+09, "floating.2.gray.hold.0.1.0.steam.entirely": -8.51427469e+08, "floating.2.gray.hold.0.1.0.steam.through": 6.95239749e+08, "floating.2.gray.hold.0.1.0.thank": true, "floating.2.gray.hold.0.1.1": "hit", "floating.2.gray.hold.0.1.2": -6.481787444899056e+08, "floating.2.gray.hold.0.1.3": 1.225027271e+09, "floating.2.gray.hold.0.1.4": -1.481507228e+09, "floating.2.gray.hold.0.1.5": true, "floating.2.gray.hold.0.2": -2.114582277e+09, "floating.2.gray.hold.0.3": 1.3900602049360588e+09, "floating.2.gray.hold.0.4": 1.6156026309049141e+09, "floating.2.gray.hold.0.5": "darkness", "floating.2.gray.hold.1": 6.3427197713988304e+07, "floating.2.gray.hold.2": -5.80344963961421e+08, "floating.2.gray.hold.3": "stems", "floating.2.gray.hold.4": 1.016960217612642e+09, "floating.2.gray.hold.5": 1.240918909e+09, "floating.2.gray.parent": "pull", "floating.2.gray.shore": -7.38396277e+08, "floating.2.gray.usually": 1.050049449e+09, "floating.2.gray.wonder": false, "floating.2.joy": "difference", "floating.2.little": "cloud", "floating.2.probably": -4.13625494e+08, "floating.2.ready": "silent", "floating.2.worker": "situation", "floating.3": "grade", "floating.4": false, "floating.5": "thou", "product": "whale", "shop": 1.294397217e+09, "spend": "greatest", "wagon": -1.722583702e+09},
		},
		{raw: []byte(`{"foo":"bar"}`), expected: map[string]interface{}{"foo": "bar"}},
		{raw: []byte(`{"foo":["bar",{"foo":"bar"}]}`), expected: map[string]interface{}{"foo.0": "bar", "foo.1.foo": "bar"}},
		{raw: []byte(`{"foo":"bar","baz":{"bar":"foo"}}`), expected: map[string]interface{}{"foo": "bar", "baz.bar": "foo"}},
		{
			raw:      []byte(`{"foo":"bar","baz":{"bar":"foo"},"bar":["foo","bar","baz"]}`),
			expected: map[string]interface{}{"bar.0": "foo", "bar.1": "bar", "bar.2": "baz", "baz.bar": "foo", "foo": "bar"},
		},
		{raw: []byte(`[]`), expected: nil},
		{raw: []byte(`null`), expected: nil},
		{raw: []byte(`"bar"`), expected: nil},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.EqualValues(t, tc.expected, Flatten(tc.raw))
		})
	}
}
