package context

import (
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/common/handler"
	"golang.org/x/net/context"
	"strings"
)

var RouterParamKey handler.Key = 0

func NewContextFromRouterParams(ctx context.Context, ps httprouter.Params) context.Context {
	return context.WithValue(ctx, RouterParamKey, ps)
}

func FetchRouterParamsFromContext(ctx context.Context, keys ...string) (map[string]string, error) {
	var r string
	res := make(map[string]string)
	ps, ok := ctx.Value(RouterParamKey).(httprouter.Params)
	if !ok {
		ps = httprouter.Params{}
	}
	for _, key := range keys {
		r = ps.ByName(key)
		if len(strings.TrimSpace(r)) == 0 {
			return map[string]string{}, errors.New(`Router param "` + key + `" empty.`)
		}
		res[key] = r
	}
	return res, nil
}
