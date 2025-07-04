// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/configx"
)

func TestContext(t *testing.T) {
	ctx := context.Background()

	actual, err := ConfigFromContext(ctx)
	require.Error(t, err)
	require.Nil(t, actual)

	assert.Panics(t, func() {
		_ = MustConfigFromContext(ctx)
	})

	expected := &configx.Provider{}
	ctx = WithConfig(ctx, expected)

	actual, err = ConfigFromContext(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	actual = MustConfigFromContext(ctx)
	require.Equal(t, expected, actual)
}

func ExampleConfigFromContext() {
	ctx := context.Background()

	config, err := configx.New(ctx, []byte(`{"type":"object","properties":{"foo":{"type":"string"}}}`), configx.WithValue("foo", "bar"))
	if err != nil {
		panic(err)
	}

	ctx = WithConfig(ctx, config)
	fmt.Printf("foo = %s", MustConfigFromContext(ctx).String("foo"))
	// Output: foo = bar
}
