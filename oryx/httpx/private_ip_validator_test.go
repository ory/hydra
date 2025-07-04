// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAssociatedIPAllowed(t *testing.T) {
	for _, disallowed := range []string{
		"localhost",
		"https://localhost/foo?bar=baz#zab",
		"127.0.0.0",
		"127.255.255.255",
		"172.16.0.0",
		"172.31.255.255",
		"192.168.0.0",
		"192.168.255.255",
		"10.0.0.0",
		"0.0.0.0",
		"10.255.255.255",
		"::1",
		"100::1",
		"fe80::1",
		"169.254.169.254", // AWS instance metadata service
	} {
		t.Run("case="+disallowed, func(t *testing.T) {
			assert.Error(t, DisallowIPPrivateAddresses(disallowed))
		})
	}
}

func TestDisallowLocalIPAddressesWhenSet(t *testing.T) {
	require.NoError(t, DisallowIPPrivateAddresses(""))
	require.Error(t, DisallowIPPrivateAddresses("127.0.0.1"))
	require.ErrorAs(t, DisallowIPPrivateAddresses("127.0.0.1"), new(ErrPrivateIPAddressDisallowed))
}

type noOpRoundTripper struct{}

func (n noOpRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

var _ http.RoundTripper = new(noOpRoundTripper)

type errRoundTripper struct{ err error }

var errNotOnWhitelist = errors.New("OK")
var errOnWhitelist = errors.New("OK (on whitelist)")

func (n errRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return nil, n.err
}

var _ http.RoundTripper = new(errRoundTripper)

// TestInternalRespectsRoundTripper tests if the RoundTripper picks the correct
// underlying transport for two allowed requests.
func TestInternalRespectsRoundTripper(t *testing.T) {
	rt := &noInternalIPRoundTripper{
		onWhitelist:    &errRoundTripper{errOnWhitelist},
		notOnWhitelist: &errRoundTripper{errNotOnWhitelist},
		internalIPExceptions: []string{
			"https://127.0.0.1/foo",
		}}

	req, err := http.NewRequest("GET", "https://google.com/foo", nil)
	require.NoError(t, err)
	_, err = rt.RoundTrip(req)
	require.ErrorIs(t, err, errNotOnWhitelist)

	req, err = http.NewRequest("GET", "https://127.0.0.1/foo", nil)
	require.NoError(t, err)
	_, err = rt.RoundTrip(req)
	require.ErrorIs(t, err, errOnWhitelist)
}

func TestAllowExceptions(t *testing.T) {
	rt := noInternalIPRoundTripper{
		onWhitelist:    &errRoundTripper{errOnWhitelist},
		notOnWhitelist: &errRoundTripper{errNotOnWhitelist},
		internalIPExceptions: []string{
			"http://localhost/asdf",
		}}

	req, err := http.NewRequest("GET", "http://localhost/asdf", nil)
	require.NoError(t, err)
	_, err = rt.RoundTrip(req)
	require.ErrorIs(t, err, errOnWhitelist)

	req, err = http.NewRequest("GET", "http://localhost/not-asdf", nil)
	require.NoError(t, err)
	_, err = rt.RoundTrip(req)
	require.ErrorIs(t, err, errNotOnWhitelist)

	req, err = http.NewRequest("GET", "http://127.0.0.1", nil)
	require.NoError(t, err)
	_, err = rt.RoundTrip(req)
	require.ErrorIs(t, err, errNotOnWhitelist)
}
