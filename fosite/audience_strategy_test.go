// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultAudienceMatchingStrategy(t *testing.T) {
	for k, tc := range []struct {
		h   []string
		n   []string
		err bool
	}{
		{
			h:   []string{},
			n:   []string{},
			err: false,
		},
		{
			h:   []string{"http://foo/bar"},
			n:   []string{},
			err: false,
		},
		{
			h:   []string{},
			n:   []string{"http://foo/bar"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users/"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users/"},
			n:   []string{"https://cloud.ory.sh/api/users/"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users/"},
			n:   []string{"https://cloud.ory.sh/api/users"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users/1234"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users", "https://cloud.ory.sh/api/users/", "https://cloud.ory.sh/api/users/1234"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users", "https://cloud.ory.sh/api/tenants"},
			n:   []string{"https://cloud.ory.sh/api/users", "https://cloud.ory.sh/api/users/", "https://cloud.ory.sh/api/users/1234", "https://cloud.ory.sh/api/tenants"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users1234"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"http://cloud.ory.sh/api/users"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh:8000/api/users"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.xyz/api/users"},
			err: true,
		},
		{
			h:   []string{"foobar"},
			n:   []string{"foobar"},
			err: false,
		},
		{
			h:   []string{"foo bar"},
			n:   []string{"foo bar"},
			err: false,
		},
		{
			h:   []string{"foobar"},
			n:   []string{"foobar"},
			err: false,
		},
		{
			h:   []string{"zoo", "bar"},
			n:   []string{"zoo"},
			err: false,
		},
		{
			h:   []string{"zoo"},
			n:   []string{"zoo", "bar"},
			err: true,
		},
		{
			h:   []string{"foobar"},
			n:   []string{"foobar/"},
			err: false,
		},
		{
			h:   []string{"foobar/"},
			n:   []string{"foobar"},
			err: false,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			err := DefaultAudienceMatchingStrategy(tc.h, tc.n)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestExactAudienceMatchingStrategy(t *testing.T) {
	for k, tc := range []struct {
		h   []string
		n   []string
		err bool
	}{
		{
			h:   []string{},
			n:   []string{},
			err: false,
		},
		{
			h:   []string{"http://foo/bar"},
			n:   []string{},
			err: false,
		},
		{
			h:   []string{},
			n:   []string{"http://foo/bar"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users/"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users/"},
			n:   []string{"https://cloud.ory.sh/api/users/"},
			err: false,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users/"},
			n:   []string{"https://cloud.ory.sh/api/users"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users/1234"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users", "https://cloud.ory.sh/api/users/", "https://cloud.ory.sh/api/users/1234"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users", "https://cloud.ory.sh/api/tenants"},
			n:   []string{"https://cloud.ory.sh/api/users", "https://cloud.ory.sh/api/users/", "https://cloud.ory.sh/api/users/1234", "https://cloud.ory.sh/api/tenants"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh/api/users1234"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"http://cloud.ory.sh/api/users"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.sh:8000/api/users"},
			err: true,
		},
		{
			h:   []string{"https://cloud.ory.sh/api/users"},
			n:   []string{"https://cloud.ory.xyz/api/users"},
			err: true,
		},
		{
			h:   []string{"foobar"},
			n:   []string{"foobar"},
			err: false,
		},
		{
			h:   []string{"foo bar"},
			n:   []string{"foo bar"},
			err: false,
		},
		{
			h:   []string{"foobar"},
			n:   []string{"foobar"},
			err: false,
		},
		{
			h:   []string{"zoo", "bar"},
			n:   []string{"zoo"},
			err: false,
		},
		{
			h:   []string{"zoo"},
			n:   []string{"zoo", "bar"},
			err: true,
		},
		{
			h:   []string{"foobar"},
			n:   []string{"foobar/"},
			err: true,
		},
		{
			h:   []string{"foobar/"},
			n:   []string{"foobar"},
			err: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			err := ExactAudienceMatchingStrategy(tc.h, tc.n)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
