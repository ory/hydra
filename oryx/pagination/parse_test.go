// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		d   string
		url string
		dl  int
		do  int
		ml  int
		el  int
		eo  int
	}{
		{"normal", "http://localhost/foo?limit=10&offset=10", 0, 0, 120, 10, 10},
		{"defaults", "http://localhost/foo", 5, 5, 10, 5, 5},
		{"defaults_and_limits", "http://localhost/foo", 5, 5, 2, 2, 5},
		{"limits", "http://localhost/foo?limit=10&offset=10", 0, 0, 5, 5, 10},
		{"negatives", "http://localhost/foo?limit=-1&offset=-1", 0, 0, 5, 0, 0},
		{"default_negatives", "http://localhost/foo", -1, -1, 5, 0, 0},
		{"invalid_defaults", "http://localhost/foo?limit=a&offset=b", 10, 10, 15, 10, 10},
	} {
		t.Run(fmt.Sprintf("case=%s", tc.d), func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			limit, offset := Parse(&http.Request{URL: u}, tc.dl, tc.do, tc.ml)
			assert.EqualValues(t, limit, tc.el)
			assert.EqualValues(t, offset, tc.eo)
		})
	}
}
