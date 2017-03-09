package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizeResponse(t *testing.T) {
	ar := NewAuthorizeResponse()
	ar.AddFragment("foo", "bar")
	ar.AddQuery("foo", "baz")
	ar.AddHeader("foo", "foo")

	ar.AddFragment("code", "bar")
	assert.Equal(t, "bar", ar.GetCode())
	ar.AddQuery("code", "baz")
	assert.Equal(t, "baz", ar.GetCode())

	assert.Equal(t, "bar", ar.GetFragment().Get("foo"))
	assert.Equal(t, "baz", ar.GetQuery().Get("foo"))
	assert.Equal(t, "foo", ar.GetHeader().Get("foo"))
}
