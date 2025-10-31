// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessRequest(t *testing.T) {
	ar := NewAccessRequest(nil)
	ar.GrantTypes = Arguments{"foobar"}
	ar.Client = &DefaultClient{}
	ar.GrantScope("foo")
	ar.SetRequestedAudience(Arguments{"foo", "foo", "bar"})
	ar.SetRequestedScopes(Arguments{"foo", "foo", "bar"})
	assert.True(t, ar.GetGrantedScopes().Has("foo"))
	assert.NotNil(t, ar.GetRequestedAt())
	assert.Equal(t, ar.GrantTypes, ar.GetGrantTypes())
	assert.Equal(t, Arguments{"foo", "bar"}, ar.RequestedAudience)
	assert.Equal(t, Arguments{"foo", "bar"}, ar.RequestedScope)
	assert.Equal(t, ar.Client, ar.GetClient())
}
