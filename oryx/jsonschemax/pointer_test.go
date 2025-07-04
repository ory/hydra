// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonschemax

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONPointerToDotNotation(t *testing.T) {
	for k, tc := range [][]string{
		{"#/foo/bar/baz", "foo.bar.baz"},
		{"#/baz", "baz"},
		{"#/properties/ory.sh~1kratos/type", "properties.ory\\.sh/kratos.type"},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			path, err := JSONPointerToDotNotation(tc[0])
			require.NoError(t, err)
			require.Equal(t, tc[1], path)
		})
	}

	_, err := JSONPointerToDotNotation("http://foo/#/bar")
	require.Error(t, err, "should fail because remote pointers are not supported")

	_, err = JSONPointerToDotNotation("http://foo/#/bar%zz")
	require.Error(t, err, "should fail because %3b is not a valid escaped path.")
}
