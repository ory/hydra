package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{}
	_ = c.Context()
}

func TestSystemSecret(t *testing.T) {
	c3 := &Config{}
	assert.EqualValues(t, c3.GetSystemSecret(), c3.GetSystemSecret())
	c := &Config{SystemSecret: "foobarbazbarasdfasdffoobarbazbarasdfasdf"}
	assert.EqualValues(t, c.GetSystemSecret(), c.GetSystemSecret())
	c2 := &Config{SystemSecret: "foobarbazbarasdfasdffoobarbazbarasdfasdf"}
	assert.EqualValues(t, c.GetSystemSecret(), c2.GetSystemSecret())
}
