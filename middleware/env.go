package middleware

import (
	"github.com/ory-am/ladon/guard/operator"
	"net/http"
	"time"
)

type Env struct {
	ctx *operator.Context
}

func NewEnv(req *http.Request) *Env {
	e := &Env{
		ctx: new(operator.Context),
	}
	e.Req(req)
	return e
}

func (e *Env) Ctx() *operator.Context {
	return e.ctx
}

func (e *Env) Req(req *http.Request) *Env {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}
	e.ctx.ClientIP = ip
	e.ctx.UserAgent = req.Header.Get("User-Agent")
	e.ctx.Timestamp = time.Now()
	return e
}

func (e *Env) Owner(owner string) *Env {
	e.ctx.Owner = owner
	return e
}
