package context

import (
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestGetterAndSetter(t *testing.T) {
	params := httprouter.Params{{"foo", "bar"}}
	ctx := NewContextFromRouterParams(context.Background(), params)
	assert.NotNil(t, ctx)
	res, err := FetchRouterParamsFromContext(ctx, "foo")
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"foo": "bar"}, res)

	_, err = FetchRouterParamsFromContext(context.Background(), "foo")
	assert.NotNil(t, err)
}
