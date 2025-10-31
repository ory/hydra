// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientWithCustomTokenLifespans(t *testing.T) {
	clc := &DefaultClientWithCustomTokenLifespans{
		DefaultClient: &DefaultClient{},
	}

	assert.Equal(t, clc.GetTokenLifespans(), (*ClientLifespanConfig)(nil))

	require.Equal(t, time.Minute*42, GetEffectiveLifespan(clc, GrantTypeImplicit, IDToken, time.Minute*42))

	customLifespan := 36 * time.Hour
	clc.SetTokenLifespans(&ClientLifespanConfig{ImplicitGrantIDTokenLifespan: &customLifespan})
	assert.NotEqual(t, clc.GetTokenLifespans(), nil)

	require.Equal(t, customLifespan, GetEffectiveLifespan(clc, GrantTypeImplicit, IDToken, time.Minute*42))
	var _ ClientWithCustomTokenLifespans = clc
}
