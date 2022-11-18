// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package spec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/jsonschema/v3"
)

func TestConfigSchema(t *testing.T) {
	c := jsonschema.NewCompiler()
	require.NoError(t, AddConfigSchema(c))

	_, err := c.Compile(context.Background(), ConfigSchemaID)
	require.NoError(t, err)
}
