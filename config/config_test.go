package config

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{}
	_ = c.Context()

	assert.Equal(t, c.GetAccessTokenLifespan(), time.Hour)
}

func TestDoesRequestSatisfyTermination(t *testing.T) {
	c := &Config{AllowTLSTermination: ""}
	assert.NotNil(t, c.DoesRequestSatisfyTermination(new(http.Request)))

	c = &Config{AllowTLSTermination: "127.0.0.1/24"}
	r := &http.Request{Header: http.Header{}}
	assert.NotNil(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{Header: http.Header{"X-Forwarded-Proto": []string{"http"}}}
	assert.NotNil(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{
		RemoteAddr: "227.0.0.1:123",
		Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
	}
	assert.NotNil(t, c.DoesRequestSatisfyTermination(r))

	r = &http.Request{
		RemoteAddr: "127.0.0.1:123",
		Header:     http.Header{"X-Forwarded-Proto": []string{"https"}},
	}
	assert.Nil(t, c.DoesRequestSatisfyTermination(r))
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
