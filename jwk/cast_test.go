// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestMustRSAPrivate(t *testing.T) {
	t.Parallel()
	keys, err := GenerateJWK(jose.RS256, "foo", "sig")
	require.NoError(t, err)

	priv := keys.Key("foo")[0]
	_, err = ToRSAPrivate(&priv)
	assert.Nil(t, err)

	MustRSAPrivate(&priv)

	pub := keys.Key("foo")[0].Public()
	_, err = ToRSAPublic(&pub)
	assert.Nil(t, err)
	MustRSAPublic(&pub)
}
