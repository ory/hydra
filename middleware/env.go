package middleware

import (
	"github.com/ory-am/ladon/guard/operator"
	"net/http"
	"time"
)

type env struct {
	ctx *operator.Context
}

func Env(req *http.Request) *env {
	return &env{
		ctx: new(operator.Context),
	}
}

func (e *env) Ctx() *operator.Context {
	return e.ctx
}

func (e *env) Req(req *http.Request) *env {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}
	e.ctx.ClientIP = ip
	e.ctx.UserAgent = req.Header.Get("User-Agent")
	e.ctx.Timestamp = time.Now()
	return e
}

func (e *env) Owner(owner string) *env {
	e.ctx.Owner = owner
	return e
}
