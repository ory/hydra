// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/driver/config"
	. "github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
)

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("should not have been called")
}

func noopHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func TestDoesRequestSatisfyTermination(t *testing.T) {
	c := testhelpers.NewConfigurationWithDefaultsAndHTTPS()
	r := testhelpers.NewRegistryMemory(t, c, &contextx.Default{})

	t.Run("case=tls-termination-disabled", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, "")

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{Header: http.Header{}, URL: new(url.URL)}, panicHandler)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	// change: x-forwarded-proto is checked after cidr, therefore it will never actually test header
	t.Run("case=missing-x-forwarded-proto", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "127.0.0.1:123",
			Header:     http.Header{},
			URL:        new(url.URL)},
			panicHandler,
		)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	// change: x-forwarded-proto is checked after cidr, therefor it will never actually test header with "http"
	t.Run("case=x-forwarded-proto-is-http", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "127.0.0.1:123",
			Header: http.Header{
				"X-Forwarded-Proto": []string{"http"},
			}, URL: new(url.URL)},
			panicHandler,
		)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	t.Run("case=missing-x-forwarded-for", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{Header: http.Header{"X-Forwarded-Proto": []string{"https"}}, URL: new(url.URL)}, panicHandler)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	t.Run("case=remote-not-in-cidr", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header:     http.Header{"X-Forwarded-Proto": []string{"https"}}, URL: new(url.URL)},
			panicHandler,
		)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	t.Run("case=remote-and-forwarded-not-in-cidr", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header: http.Header{
				"X-Forwarded-Proto": []string{"https"},
				"X-Forwarded-For":   []string{"227.0.0.1"},
			}, URL: new(url.URL)},
			panicHandler,
		)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})

	t.Run("case=remote-matches-cidr", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "127.0.0.1:123",
			Header: http.Header{
				"X-Forwarded-Proto": []string{"https"},
			}, URL: new(url.URL)},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	// change: cidr and x-forwarded-proto headers are irrelevant for this test
	t.Run("case=passes-because-health-alive-endpoint", func(t *testing.T) {
		c.MustSet(context.Background(), config.AdminInterface.Key(config.KeySuffixTLSAllowTerminationFrom), []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.AdminInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header:     http.Header{},
			URL:        &url.URL{Path: "/health/alive"},
		},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	// change: cidr and x-forwarded-proto headers are irrelevant for this test
	t.Run("case=passes-because-health-ready-endpoint", func(t *testing.T) {
		c.MustSet(context.Background(), config.AdminInterface.Key(config.KeySuffixTLSAllowTerminationFrom), []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.AdminInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header:     http.Header{},
			URL:        &url.URL{Path: "/health/alive"},
		},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	t.Run("case=forwarded-matches-cidr", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.2:123",
			Header: http.Header{
				"X-Forwarded-For":   []string{"227.0.0.1, 127.0.0.1, 227.0.0.2"},
				"X-Forwarded-Proto": []string{"https"},
			}, URL: new(url.URL)},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	t.Run("case=forwarded-matches-cidr-without-spaces", func(t *testing.T) {
		c.MustSet(context.Background(), config.KeyTLSAllowTerminationFrom, []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.2:123",
			Header: http.Header{
				"X-Forwarded-For":   []string{"227.0.0.1,127.0.0.1,227.0.0.2"},
				"X-Forwarded-Proto": []string{"https"},
			}, URL: new(url.URL)},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	// test: in case http is forced request should be accepted
	t.Run("case=forced-http", func(t *testing.T) {
		c := testhelpers.NewConfigurationWithDefaults()
		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.PublicInterface))(res, &http.Request{Header: http.Header{}, URL: new(url.URL)}, noopHandler)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	// test: prometheus endpoint should accept request
	t.Run("case=passes-with-tls-upstream-on-metrics-prometheus-endpoint", func(t *testing.T) {
		c.MustSet(context.Background(), config.AdminInterface.Key(config.KeySuffixTLSAllowTerminationFrom), []string{"126.0.0.1/24", "127.0.0.1/24"})

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.AdminInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header:     http.Header{},
			URL:        &url.URL{Path: "/metrics/prometheus"},
		},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	// test: prometheus endpoint should accept request because TLS is disabled
	t.Run("case=passes-with-tls-disabled-on-admin-endpoint", func(t *testing.T) {
		c.MustSet(context.Background(), config.AdminInterface.Key(config.KeySuffixTLSEnabled), false)

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.AdminInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header: http.Header{
				"X-Forwarded-Proto": []string{"http"},
			},
			URL: &url.URL{Path: "/foo"},
		},
			noopHandler,
		)
		assert.EqualValues(t, http.StatusNoContent, res.Code)
	})

	// test: prometheus endpoint should not accept request because TLS is enabled
	t.Run("case=fails-with-tls-enabled-on-admin-endpoint", func(t *testing.T) {
		c.MustSet(context.Background(), config.AdminInterface.Key(config.KeySuffixTLSEnabled), true)

		res := httptest.NewRecorder()
		RejectInsecureRequests(r, c.TLS(context.Background(), config.AdminInterface))(res, &http.Request{
			RemoteAddr: "227.0.0.1:123",
			Header: http.Header{
				"X-Forwarded-Proto": []string{"http"},
			},
			URL: &url.URL{Path: "/foo"},
		},
			panicHandler,
		)
		assert.EqualValues(t, http.StatusBadGateway, res.Code)
	})
}
