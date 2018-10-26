/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package config

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{
		DatabaseURL: "memory",
	}
	_ = c.Context()

	assert.Equal(t, c.GetAccessTokenLifespan(), time.Hour)
}

func TestDoesRequestSatisfyTermination(t *testing.T) {
	c := &Config{AllowTLSTermination: ""}
	assert.Error(t, c.DoesRequestSatisfyTermination(&http.Request{Header: http.Header{}, URL: new(url.URL)}))

	c = &Config{AllowTLSTermination: "127.0.0.1/24"}
	r := &http.Request{Header: http.Header{}, URL: new(url.URL)}
	assert.Error(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{Header: http.Header{"X-Forwarded-Proto": []string{"http"}}, URL: new(url.URL)}
	assert.Error(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{
		RemoteAddr: "227.0.0.1:123",
		Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
		URL:        new(url.URL),
	}
	assert.Error(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{
		RemoteAddr: "127.0.0.1:123",
		Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
		URL:        new(url.URL),
	}
	assert.NoError(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{
		RemoteAddr: "127.0.0.1:123",
		Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
		URL:        &url.URL{Path: "/health"},
	}
	assert.NoError(t, c.DoesRequestSatisfyTermination(r))
}

func TestTracingSetup(t *testing.T) {
	// tracer is not loaded if an unknown tracing provider is specified
	c := &Config{TracingProvider: "some_unsupported_tracing_provider"}
	tracer, _ := c.GetTracer()
	assert.False(t, tracer.IsLoaded())
	assert.False(t, c.WithTracing())

	// tracer is not loaded if no tracing provider is specified
	c = &Config{TracingProvider: ""}
	tracer, _ = c.GetTracer()
	assert.False(t, tracer.IsLoaded())
	assert.False(t, c.WithTracing())

	// tracer is loaded if configured properly
	c = &Config{
		TracingProvider:          "jaeger",
		TracingServiceName:       "Ory Hydra",
		JaegerSamplingServerUrl:  "http://localhost:5778/sampling",
		JaegerLocalAgentHostPort: "127.0.0.1:6831",
	}
	tracer, _ = c.GetTracer()
	assert.True(t, tracer.IsLoaded())
	assert.True(t, c.WithTracing())
}

func TestSystemSecret(t *testing.T) {
	c3 := &Config{}
	assert.EqualValues(t, c3.GetSystemSecret(), c3.GetSystemSecret())
	c := &Config{SystemSecret: "foobarbazbarasdfasdffoobarbazbarasdfasdf"}
	assert.EqualValues(t, c.GetSystemSecret(), c.GetSystemSecret())
	c2 := &Config{SystemSecret: "foobarbazbarasdfasdffoobarbazbarasdfasdf"}
	assert.EqualValues(t, c.GetSystemSecret(), c2.GetSystemSecret())
}

func TestResolve(t *testing.T) {
	c := &Config{EndpointURL: "https://localhost:1234"}
	assert.Equal(t, c.Resolve("foo", "bar").String(), "https://localhost:1234/foo/bar")
	assert.Equal(t, c.Resolve("/foo", "/bar").String(), "https://localhost:1234/foo/bar")

	c = &Config{EndpointURL: "https://localhost:1234/"}
	assert.Equal(t, c.Resolve("/foo", "/bar").String(), "https://localhost:1234/foo/bar")

	c = &Config{EndpointURL: "https://localhost:1234/bar"}
	assert.Equal(t, c.Resolve("/foo", "/bar").String(), "https://localhost:1234/bar/foo/bar")
}

func TestLifespan(t *testing.T) {
	assert.Equal(t, (&Config{}).GetAccessTokenLifespan(), time.Hour)
	assert.Equal(t, (&Config{AccessTokenLifespan: "6h"}).GetAccessTokenLifespan(), time.Hour*6)

	assert.Equal(t, (&Config{}).GetAuthCodeLifespan(), time.Minute*10)
	assert.Equal(t, (&Config{AuthCodeLifespan: "15m"}).GetAuthCodeLifespan(), time.Minute*15)

	assert.Equal(t, (&Config{}).GetIDTokenLifespan(), time.Hour)
	assert.Equal(t, (&Config{IDTokenLifespan: "10s"}).GetIDTokenLifespan(), time.Second*10)
}
