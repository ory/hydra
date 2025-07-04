// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package logrusx

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/sjson"

	"github.com/ory/jsonschema/v3"
)

func TestConfigSchema(t *testing.T) {
	config := func(t *testing.T, vals map[string]interface{}) []byte {
		rawConfig, err := sjson.Set("{}", "log", vals)
		require.NoError(t, err)

		return []byte(rawConfig)
	}

	t.Run("case=basic validation and retrieval", func(t *testing.T) {
		c := jsonschema.NewCompiler()
		require.NoError(t, AddConfigSchema(c))
		schema, err := c.Compile(context.Background(), ConfigSchemaID)
		require.NoError(t, err)

		logConfig := map[string]interface{}{
			"level":                 "trace",
			"format":                "json_pretty",
			"leak_sensitive_values": true,
			"additional_redacted_headers": []interface{}{
				"custom_header_1",
				"custom_header_2",
			},
		}
		assert.NoError(t, schema.ValidateInterface(logConfig))

		k := koanf.New(".")
		require.NoError(t, k.Load(rawbytes.Provider(config(t, logConfig)), json.Parser()))

		l := New("foo", "bar", WithConfigurator(k))

		assert.True(t, l.leakSensitive)
		assert.Equal(t, logrus.TraceLevel, l.Logger.Level)
		assert.Contains(t, l.additionalRedactedHeaders, "custom_header_1")
		assert.Contains(t, l.additionalRedactedHeaders, "custom_header_2")
		assert.IsType(t, &logrus.JSONFormatter{}, l.Logger.Formatter)
	})

	t.Run("case=warns on unknown format", func(t *testing.T) {
		h := &test.Hook{}
		New("foo", "bar", WithHook(h), ForceFormat("unknown"))

		require.Len(t, h.Entries, 1)
		assert.Contains(t, h.LastEntry().Message, "got unknown \"log.format\", falling back to \"text\"")
	})

	t.Run("case=does not warn on text format", func(t *testing.T) {
		h := &test.Hook{}
		New("foo", "bar", WithHook(h), ForceFormat("text"))

		assert.Len(t, h.Entries, 0)
	})
}
