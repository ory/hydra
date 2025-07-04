// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/inhies/go-bytesize"

	"github.com/knadh/koanf/parsers/json"

	"github.com/ory/x/urlx"

	"github.com/spf13/pflag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newProvider(t testing.TB) *Provider {
	// Fake some flags
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.String("foo-bar-baz", "", "")
	f.StringP("b", "b", "", "")
	args := []string{"/var/folders/mt/m1dwr59n73zgsq7bk0q2lrmc0000gn/T/go-build533083141/b001/exe/asdf", "aaaa", "-b", "bbbb", "dddd", "eeee", "--foo-bar-baz", "fff"}
	require.NoError(t, f.Parse(args[1:]))
	RegisterFlags(f)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	p, err := New(ctx, []byte(`{"type": "object", "properties": {"foo-bar-baz": {"type": "string"}, "b": {"type": "string"}}}`), WithFlags(f), WithContext(ctx))
	require.NoError(t, err)
	return p
}

func TestProviderMethods(t *testing.T) {
	p := newProvider(t)

	t.Run("check flags", func(t *testing.T) {
		assert.Equal(t, "fff", p.String("foo-bar-baz"))
		assert.Equal(t, "bbbb", p.String("b"))
	})

	t.Run("check fallbacks", func(t *testing.T) {
		t.Run("type=string", func(t *testing.T) {
			require.NoError(t, p.Set("some.string", "bar"))
			assert.Equal(t, "bar", p.StringF("some.string", "baz"))
			assert.Equal(t, "baz", p.StringF("not.some.string", "baz"))
		})

		t.Run("type=float", func(t *testing.T) {
			require.NoError(t, p.Set("some.float", 123.123))
			assert.Equal(t, 123.123, p.Float64F("some.float", 321.321))
			assert.Equal(t, 321.321, p.Float64F("not.some.float", 321.321))
		})

		t.Run("type=int", func(t *testing.T) {
			require.NoError(t, p.Set("some.int", 123))
			assert.Equal(t, 123, p.IntF("some.int", 123))
			assert.Equal(t, 321, p.IntF("not.some.int", 321))
		})

		t.Run("type=bytesize", func(t *testing.T) {
			const key = "some.bytesize"

			for _, v := range []interface{}{
				bytesize.MB,
				float64(1024 * 1024),
				"1MB",
			} {
				require.NoError(t, p.Set(key, v))
				assert.Equal(t, bytesize.MB, p.ByteSizeF(key, 0))
			}
		})

		github := urlx.ParseOrPanic("https://github.com/ory")
		ory := urlx.ParseOrPanic("https://www.ory.sh/")

		t.Run("type=url", func(t *testing.T) {
			require.NoError(t, p.Set("some.url", "https://github.com/ory"))
			assert.Equal(t, github, p.URIF("some.url", ory))
			assert.Equal(t, ory, p.URIF("not.some.url", ory))
		})

		t.Run("type=request_uri", func(t *testing.T) {
			require.NoError(t, p.Set("some.request_uri", "https://github.com/ory"))
			assert.Equal(t, github, p.RequestURIF("some.request_uri", ory))
			assert.Equal(t, ory, p.RequestURIF("not.some.request_uri", ory))

			require.NoError(t, p.Set("invalid.request_uri", "foo"))
			assert.Equal(t, ory, p.RequestURIF("invalid.request_uri", ory))
		})
	})

	t.Run("allow integer as duration", func(t *testing.T) {
		assert.NoError(t, p.Set("duration.integer1", -1))
		assert.NoError(t, p.Set("duration.integer2", "-1"))

		assert.Equal(t, -1*time.Nanosecond, p.DurationF("duration.integer1", time.Second))
		assert.Equal(t, -1*time.Nanosecond, p.DurationF("duration.integer2", time.Second))
	})

	t.Run("use complex set operations", func(t *testing.T) {
		assert.NoError(t, p.Set("nested", nil))
		assert.NoError(t, p.Set("nested.value", "https://www.ory.sh/kratos"))
		assert.Equal(t, "https://www.ory.sh/kratos", p.Get("nested.value"))
	})

	t.Run("use DirtyPatch operations", func(t *testing.T) {
		assert.NoError(t, p.DirtyPatch("nested", nil))
		assert.NoError(t, p.DirtyPatch("nested.value", "https://www.ory.sh/kratos"))
		assert.Equal(t, "https://www.ory.sh/kratos", p.Get("nested.value"))

		assert.NoError(t, p.DirtyPatch("duration.integer1", -1))
		assert.NoError(t, p.DirtyPatch("duration.integer2", "-1"))
		assert.Equal(t, -1*time.Nanosecond, p.DurationF("duration.integer1", time.Second))
		assert.Equal(t, -1*time.Nanosecond, p.DurationF("duration.integer2", time.Second))

		require.NoError(t, p.DirtyPatch("some.float", 123.123))
		assert.Equal(t, 123.123, p.Float64F("some.float", 321.321))
		assert.Equal(t, 321.321, p.Float64F("not.some.float", 321.321))
	})
}

func TestAdvancedConfigs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, tc := range []struct {
		stub      string
		configs   []string
		envs      [][2]string
		ops       []OptionModifier
		isValid   bool
		expectedF func(*testing.T, *Provider)
	}{
		{
			stub:    "nested-array",
			configs: []string{"stub/nested-array/kratos.yaml"},
			isValid: true, envs: [][2]string{
				{"PROVIDERS_0_CLIENT_ID", "client@example.com"},
				{"PROVIDERS_1_CLIENT_ID", "some@example.com"},
			},
		},
		{
			stub:    "kratos",
			configs: []string{"stub/kratos/kratos.yaml"},
			isValid: true, envs: [][2]string{
				{"SELFSERVICE_METHODS_OIDC_CONFIG_PROVIDERS", `[{"id":"google","provider":"google","mapper_url":"file:///etc/config/kratos/oidc.google.jsonnet","client_id":"client@example.com","client_secret":"secret"}]`},
				{"DSN", "sqlite:///var/lib/sqlite/db.sqlite?_fk=true"},
				{"SELFSERVICE_FLOWS_REGISTRATION_AFTER_PASSWORD_HOOKS_0_HOOK", "session"},
			},
		},
		{
			stub:    "multi",
			configs: []string{"stub/multi/a.yaml", "stub/multi/b.yaml"},
			isValid: true, envs: [][2]string{
				{"DSN", "sqlite:///var/lib/sqlite/db.sqlite?_fk=true"},
			}},
		{
			stub:    "from-files",
			isValid: true, envs: [][2]string{
				{"DSN", "sqlite:///var/lib/sqlite/db.sqlite?_fk=true"},
			},
			ops: []OptionModifier{WithConfigFiles("stub/multi/a.yaml", "stub/multi/b.yaml")}},
		{
			stub:    "hydra",
			configs: []string{"stub/hydra/hydra.yaml"},
			isValid: true,
			envs: [][2]string{
				{"DSN", "sqlite:///var/lib/sqlite/db.sqlite?_fk=true"},
				{"TRACING_PROVIDER", "jaeger"},
				{"TRACING_PROVIDERS_JAEGER_SAMPLING_SERVER_URL", "http://jaeger:5778/sampling"},
				{"TRACING_PROVIDERS_JAEGER_LOCAL_AGENT_ADDRESS", "jaeger:6831"},
				{"TRACING_PROVIDERS_JAEGER_SAMPLING_TYPE", "const"},
				{"TRACING_PROVIDERS_JAEGER_SAMPLING_VALUE", "1"},
			},
			expectedF: func(t *testing.T, p *Provider) {
				assert.Equal(t, "sqlite:///var/lib/sqlite/db.sqlite?_fk=true", p.Get("dsn"))
				assert.Equal(t, "jaeger", p.Get("tracing.provider"))
			}},
		{
			stub:    "hydra",
			configs: []string{"stub/hydra/hydra.yaml"},
			isValid: false,
			ops:     []OptionModifier{WithUserProviders(NewKoanfMemory(ctx, []byte(`{"dsn": null}`)))},
		},
		{
			stub:    "hydra",
			configs: []string{"stub/hydra/hydra.yaml"},
			isValid: true,
			ops:     []OptionModifier{WithUserProviders(NewKoanfMemory(ctx, []byte(`{"dsn": "invalid"}`)))},
			envs: [][2]string{
				{"DSN", "sqlite:///var/lib/sqlite/db.sqlite?_fk=true"},
				{"TRACING_PROVIDER", "jaeger"},
				{"TRACING_PROVIDERS_JAEGER_LOCAL_AGENT_ADDRESS", "jaeger:6831"},
				{"TRACING_PROVIDERS_JAEGER_SAMPLING_SERVER_URL", "http://jaeger:5778/sampling"},
				{"TRACING_PROVIDERS_JAEGER_SAMPLING_TYPE", "const"},
				{"TRACING_PROVIDERS_JAEGER_SAMPLING_VALUE", "1"},
			},
		},
	} {
		t.Run("service="+tc.stub, func(t *testing.T) {
			setEnvs(t, tc.envs)

			expected, err := os.ReadFile(path.Join("stub", tc.stub, "expected.json"))
			require.NoError(t, err)

			schemaPath := path.Join("stub", tc.stub, "config.schema.json")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			k, err := newKoanf(ctx, schemaPath, tc.configs, append(tc.ops, WithContext(ctx))...)
			if !tc.isValid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			out, err := k.Koanf.Marshal(json.Parser())
			require.NoError(t, err)
			assert.JSONEq(t, string(expected), string(out), "%s", out)

			if tc.expectedF != nil {
				tc.expectedF(t, k)
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	// Benchmark set function
	p := newProvider(b)
	var err error
	for i := 0; i < b.N; i++ {
		err = p.Set("nested.value", "https://www.ory.sh/kratos")
		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
	}
}

func BenchmarkDirtyPatch(b *testing.B) {
	// Benchmark set function
	p := newProvider(b)
	var err error
	for i := 0; i < b.N; i++ {
		err = p.DirtyPatch("nested.value", "https://www.ory.sh/kratos")
		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
	}
}
