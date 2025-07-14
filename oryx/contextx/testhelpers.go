// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/gofrs/uuid"

	"github.com/ory/x/configx"
)

type (
	TestConfigProvider struct {
		ConfigSchema []byte
		Options      []configx.OptionModifier
	}
	contextKey int
)

func NewTestConfigProvider(schema []byte, opts ...configx.OptionModifier) *TestConfigProvider {
	return &TestConfigProvider{
		ConfigSchema: schema,
		Options:      opts,
	}
}

func (t *TestConfigProvider) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	return (&Default{}).Network(ctx, network)
}

func (t *TestConfigProvider) Config(ctx context.Context, config *configx.Provider) *configx.Provider {
	values, ok := ctx.Value(contextConfigKey).([]map[string]any)
	if !ok {
		return config
	}

	opts := make([]configx.OptionModifier, 1, 1+len(values))
	opts[0] = configx.WithValues(config.All())
	for _, v := range values {
		opts = append(opts, configx.WithValues(v))
	}
	config, err := configx.New(ctx, t.ConfigSchema, append(t.Options, opts...)...)
	if err != nil {
		// This is not production code. The provider is only used in tests.
		panic(err)
	}
	return config
}

const contextConfigKey contextKey = 1

var (
	_ Contextualizer = (*TestConfigProvider)(nil)
)

func WithConfigValue(ctx context.Context, key string, value any) context.Context {
	return WithConfigValues(ctx, map[string]any{key: value})
}

func WithConfigValues(ctx context.Context, setValues ...map[string]any) context.Context {
	values, ok := ctx.Value(contextConfigKey).([]map[string]any)
	if !ok {
		values = make([]map[string]any, 0)
	}
	newValues := make([]map[string]any, len(values), len(values)+len(setValues))
	copy(newValues, values)
	newValues = append(newValues, setValues...)

	return context.WithValue(ctx, contextConfigKey, newValues)
}

type ConfigurableTestHandler struct {
	configs map[uuid.UUID][]map[string]any
	handler http.Handler
}

func NewConfigurableTestHandler(h http.Handler) *ConfigurableTestHandler {
	return &ConfigurableTestHandler{
		configs: make(map[uuid.UUID][]map[string]any),
		handler: h,
	}
}

func (t *ConfigurableTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cID := r.Header.Get("Test-Config-Id")
	if config, ok := t.configs[uuid.FromStringOrNil(cID)]; ok {
		r = r.WithContext(WithConfigValues(r.Context(), config...))
	}
	t.handler.ServeHTTP(w, r)
}

func (t *ConfigurableTestHandler) RegisterConfig(config ...map[string]any) uuid.UUID {
	id := uuid.Must(uuid.NewV4())
	t.configs[id] = config
	return id
}

func (t *ConfigurableTestHandler) UseConfig(r *http.Request, id uuid.UUID) *http.Request {
	r.Header.Set("Test-Config-Id", id.String())
	return r
}

func (t *ConfigurableTestHandler) UseConfigValues(r *http.Request, values ...map[string]any) *http.Request {
	return t.UseConfig(r, t.RegisterConfig(values...))
}

type ConfigurableTestServer struct {
	*httptest.Server
	handler   *ConfigurableTestHandler
	transport http.RoundTripper
}

func NewConfigurableTestServer(h http.Handler) *ConfigurableTestServer {
	handler := NewConfigurableTestHandler(h)
	server := httptest.NewServer(handler)

	t := server.Client().Transport
	cts := &ConfigurableTestServer{
		handler:   handler,
		Server:    server,
		transport: t,
	}
	server.Client().Transport = cts
	return cts
}

func (t *ConfigurableTestServer) RoundTrip(r *http.Request) (*http.Response, error) {
	config, ok := r.Context().Value(contextConfigKey).([]map[string]any)
	if ok && config != nil {
		r = t.handler.UseConfigValues(r, config...)
	}
	return t.transport.RoundTrip(r)
}

type AutoContextClient struct {
	*http.Client
	transport http.RoundTripper
	ctx       context.Context
}

func (t *ConfigurableTestServer) Client(ctx context.Context) *AutoContextClient {
	baseClient := *t.Server.Client()
	autoClient := &AutoContextClient{
		Client:    &baseClient,
		transport: t,
		ctx:       ctx,
	}
	baseClient.Transport = autoClient
	return autoClient
}

func (c *AutoContextClient) RoundTrip(r *http.Request) (*http.Response, error) {
	return c.transport.RoundTrip(r.WithContext(c.ctx))
}
