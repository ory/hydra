package health

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/x/healthx"
)

type AliveCheck struct {
	Handler *healthx.Handler
}

type ReadyCheck struct {
	Handler    *healthx.Handler
	ShowErrors bool
}

type HealthCheck struct {
	Alive *AliveCheck
	Ready *ReadyCheck
}

func NewAliveCheck(handler *healthx.Handler) *AliveCheck {
	return &AliveCheck{
		Handler: handler,
	}
}

func NewReadyCheck(handler *healthx.Handler, showErrors bool) *ReadyCheck {
	return &ReadyCheck{
		ShowErrors: showErrors,
		Handler:    handler,
	}
}

func (alive AliveCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	alive.Handler.Alive(w, r, params)
}

func (ready ReadyCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	ready.Handler.Ready(ready.ShowErrors)(w, r, params)
}

func (h *HealthCheck) RegisterHealthHandlers(handler *healthx.Handler, showErrors bool) {
	h.Ready = NewReadyCheck(handler, showErrors)
	h.Alive = NewAliveCheck(handler)
}
