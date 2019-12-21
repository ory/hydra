package config

import (
	"net/http"
	"testing"
	"time"

	"net/url"

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

func TestSystemSecret(t *testing.T) {
	c3 := &Config{}
	assert.EqualValues(t, c3.GetSystemSecret(), c3.GetSystemSecret())
	c := &Config{SystemSecret: "foobarbazbarasdfasdffoobarbazbarasdfasdf"}
	assert.EqualValues(t, c.GetSystemSecret(), c.GetSystemSecret())
	c2 := &Config{SystemSecret: "foobarbazbarasdfasdffoobarbazbarasdfasdf"}
	assert.EqualValues(t, c.GetSystemSecret(), c2.GetSystemSecret())
}

func TestResolve(t *testing.T) {
	c := &Config{ClusterURL: "https://localhost:1234"}
	assert.Equal(t, c.Resolve("foo", "bar").String(), "https://localhost:1234/foo/bar")
	assert.Equal(t, c.Resolve("/foo", "/bar").String(), "https://localhost:1234/foo/bar")

	c = &Config{ClusterURL: "https://localhost:1234/"}
	assert.Equal(t, c.Resolve("/foo", "/bar").String(), "https://localhost:1234/foo/bar")

	c = &Config{ClusterURL: "https://localhost:1234/bar"}
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

func TestExtractPortFromDSN(t *testing.T) {
	m, p := extractPortFromDSN("")
	assert.Equal(t, m, "")
	assert.Equal(t, p, "")

	m, p = extractPortFromDSN("mysql://u:p@tcp(hostname:port)/db?parseTime=true")
	assert.Equal(t, m, "mysql://u:p@tcp(hostname)/db?parseTime=true")
	assert.Equal(t, p, "port")

	m, p = extractPortFromDSN("mysql://u:p@tcp(hostname)/db?parseTime=true")
	assert.Equal(t, m, "mysql://u:p@tcp(hostname)/db?parseTime=true")
	assert.Equal(t, p, "")
}
