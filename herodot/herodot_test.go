package herodot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	assert.NotNil(t, ctx.Value(RequestIDKey))

	ctx = Context(context.Background())
	assert.NotNil(t, ctx.Value(RequestIDKey))
}
