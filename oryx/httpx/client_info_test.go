// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIgnoresInternalIPs(t *testing.T) {
	input := "54.155.246.232,10.145.1.10"

	res, err := GetClientIPAddressesWithoutInternalIPs(strings.Split(input, ","))
	require.NoError(t, err)
	assert.Equal(t, "54.155.246.232", res)
}

func TestEmptyInputArray(t *testing.T) {
	res, err := GetClientIPAddressesWithoutInternalIPs([]string{})
	require.NoError(t, err)
	assert.Equal(t, "", res)
}

func TestClientIP(t *testing.T) {
	req := http.Request{
		RemoteAddr: "1.0.0.4",
		Header:     http.Header{},
	}
	req.Header.Add("true-client-ip", "1.0.0.1")
	req.Header.Add("cf-connecting-ip", "1.0.0.2")
	req.Header.Add("x-real-ip", "1.0.0.3")
	req.Header.Add("x-forwarded-for", "192.168.1.1,1.0.0.3,10.0.0.1")
	t.Run("true-client-ip", func(t *testing.T) {
		req := req.Clone(context.Background())
		assert.Equal(t, "1.0.0.1", ClientIP(req))
	})
	t.Run("cf-connecting-ip", func(t *testing.T) {
		req := req.Clone(context.Background())
		req.Header.Del("true-client-ip")
		assert.Equal(t, "1.0.0.2", ClientIP(req))
	})
	t.Run("x-real-ip", func(t *testing.T) {
		req := req.Clone(context.Background())
		req.Header.Del("true-client-ip")
		req.Header.Del("cf-connecting-ip")
		assert.Equal(t, "1.0.0.3", ClientIP(req))
	})
	t.Run("x-forwarded-for", func(t *testing.T) {
		req := req.Clone(context.Background())
		req.Header.Del("true-client-ip")
		req.Header.Del("cf-connecting-ip")
		req.Header.Del("x-real-ip")
		assert.Equal(t, "1.0.0.3", ClientIP(req))
	})
	t.Run("remote-addr", func(t *testing.T) {
		req := req.Clone(context.Background())
		req.Header.Del("true-client-ip")
		req.Header.Del("cf-connecting-ip")
		req.Header.Del("x-real-ip")
		req.Header.Del("x-forwarded-for")
		assert.Equal(t, "1.0.0.4", ClientIP(req))
	})
}

func TestClientGeoLocation(t *testing.T) {
	req := http.Request{
		Header: http.Header{},
	}
	req.Header.Add("cf-ipcity", "Berlin")
	req.Header.Add("cf-ipcountry", "Germany")
	req.Header.Add("cf-region-code", "BE")

	t.Run("cf-ipcity", func(t *testing.T) {
		req := req.Clone(context.Background())
		assert.Equal(t, "Berlin", ClientGeoLocation(req).City)
	})

	t.Run("cf-ipcountry", func(t *testing.T) {
		req := req.Clone(context.Background())
		assert.Equal(t, "Germany", ClientGeoLocation(req).Country)
	})

	t.Run("cf-region-code", func(t *testing.T) {
		req := req.Clone(context.Background())
		assert.Equal(t, "BE", ClientGeoLocation(req).Region)
	})

	t.Run("empty", func(t *testing.T) {
		req := req.Clone(context.Background())
		req.Header.Del("cf-ipcity")
		req.Header.Del("cf-ipcountry")
		req.Header.Del("cf-region-code")
		assert.Equal(t, GeoLocation{}, *ClientGeoLocation(req))
	})
}
