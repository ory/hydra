package path

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	assert.Equal(t, "http://foo/bar/baz", Join("http://foo/", "/bar", "/baz"))
	assert.Equal(t, "http://foo/bar/baz", Join("http://foo", "bar", "/baz"))
}