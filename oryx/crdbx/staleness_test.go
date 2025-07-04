// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package crdbx

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/x/urlx"
)

func TestConsistencyLevelFromString(t *testing.T) {
	assert.Equal(t, ConsistencyLevelUnset, ConsistencyLevelFromString(""))
	assert.Equal(t, ConsistencyLevelStrong, ConsistencyLevelFromString("strong"))
	assert.Equal(t, ConsistencyLevelEventual, ConsistencyLevelFromString("eventual"))
	assert.Equal(t, ConsistencyLevelStrong, ConsistencyLevelFromString("lol"))
}

func TestConsistencyLevelFromRequest(t *testing.T) {
	assert.Equal(t, ConsistencyLevelStrong, ConsistencyLevelFromRequest(&http.Request{URL: urlx.ParseOrPanic("/?consistency=strong")}))
	assert.Equal(t, ConsistencyLevelEventual, ConsistencyLevelFromRequest(&http.Request{URL: urlx.ParseOrPanic("/?consistency=eventual")}))
	assert.Equal(t, ConsistencyLevelStrong, ConsistencyLevelFromRequest(&http.Request{URL: urlx.ParseOrPanic("/?consistency=asdf")}))
	assert.Equal(t, ConsistencyLevelUnset, ConsistencyLevelFromRequest(&http.Request{URL: urlx.ParseOrPanic("/?consistency")}))

}

func TestGetTransactionConsistency(t *testing.T) {
	for k, tc := range []struct {
		in       ConsistencyLevel
		fallback ConsistencyLevel
		dialect  string
		expected string
	}{
		{
			in:       ConsistencyLevelUnset,
			fallback: ConsistencyLevelStrong,
			dialect:  "cockroach",
			expected: "",
		},
		{
			in:       ConsistencyLevelStrong,
			fallback: ConsistencyLevelStrong,
			dialect:  "cockroach",
			expected: "",
		},
		{
			in:       ConsistencyLevelStrong,
			fallback: ConsistencyLevelEventual,
			dialect:  "cockroach",
			expected: "",
		},
		{
			in:       ConsistencyLevelUnset,
			fallback: ConsistencyLevelEventual,
			dialect:  "cockroach",
			expected: transactionFollowerReadTimestamp,
		},
		{
			in:       ConsistencyLevelEventual,
			fallback: ConsistencyLevelEventual,
			dialect:  "cockroach",
			expected: transactionFollowerReadTimestamp,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			q := getTransactionConsistencyQuery(tc.dialect, tc.in, tc.fallback)
			assert.EqualValues(t, tc.expected, q)
		})
	}
}
