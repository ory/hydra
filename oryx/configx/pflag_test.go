// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"context"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/jsonschema/v3"
)

func TestPFlagProvider(t *testing.T) {
	const schema = `
{
  "type": "object",
  "properties": {
	"foo": {
	  "type": "string"
	}
  }
}
`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := jsonschema.CompileString(ctx, "", schema)
	require.NoError(t, err)

	t.Run("only parses known flags", func(t *testing.T) {
		flags := pflag.NewFlagSet("", pflag.ContinueOnError)
		flags.String("foo", "", "")
		flags.String("bar", "", "")
		require.NoError(t, flags.Parse([]string{"--foo", "x", "--bar", "y"}))

		p, err := NewPFlagProvider([]byte(schema), s, flags, nil)
		require.NoError(t, err)

		values, err := p.Read()
		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"foo": "x",
		}, values)
	})
}
