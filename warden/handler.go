package warden

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/ladon"
)

const (
	TokenValidHandlerPath   = "/warden/token/valid"
	TokenAllowedHandlerPath = "/warden/token/allowed"
	AllowedHandlerPath      = "/warden/allowed"
	IntrospectPath          = "/oauth2/introspect"
)

type WardenHandler struct {
	H      herodot.Herodot
	Warden firewall.Firewall
}

func NewHandler(c *config.Config, router *httprouter.Router) *WardenHandler {
	ctx := c.Context()

	h := &WardenHandler{
		H:      &herodot.JSON{},
		Warden: ctx.Warden,
	}
	h.SetRoutes(router)

	return h
}

type WardenAuthorizedRequest struct {
	Scopes []string `json:"scopes"`
	Token  string   `json:"token"`
}

type WardenAccessRequest struct {
	*ladon.Request
	*WardenAuthorizedRequest
}

var notAllowed = struct {
	Allowed bool `json:"allowed"`
}{Allowed: false}

var invalid = struct {
	Valid bool `json:"valid"`
}{Valid: false}

var inactive = struct {
	Active bool `json:"active"`
}{Active: false}

func (h *WardenHandler) SetRoutes(r *httprouter.Router) {
	r.POST(TokenValidHandlerPath, h.TokenValid)
	r.POST(TokenAllowedHandlerPath, h.TokenAllowed)
	r.POST(AllowedHandlerPath, h.Allowed)
	r.POST(IntrospectPath, h.Introspect)
}

func (h *WardenHandler) Introspect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := herodot.NewContext()
	clientCtx, err := h.Warden.InspectToken(ctx, TokenFromRequest(r))
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	var ar struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}
	defer r.Body.Close()

	auth, err := h.Warden.IntrospectToken(ctx, ar.Token)
	if err != nil {
		h.H.Write(ctx, w, r, &inactive)
		return
	} else if clientCtx.Subject != auth.Audience {
		h.H.Write(ctx, w, r, &inactive)
		return
	}

	h.H.Write(ctx, w, r, auth)
}

func (h *WardenHandler) TokenValid(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := herodot.NewContext()
	clientCtx, err := h.Warden.TokenAllowed(ctx, h.Warden.TokenFromRequest(r), &ladon.Request{
		Resource: "rn:hydra:warden:token:valid",
		Action:   "decide",
	}, "hydra.warden")
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	var ar WardenAuthorizedRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}
	defer r.Body.Close()

	authContext, err := h.Warden.InspectToken(ctx, ar.Token, ar.Scopes...)
	if err != nil {
		h.H.Write(ctx, w, r, &invalid)
		return
	}

	authContext.Audience = clientCtx.Subject
	h.H.Write(ctx, w, r, struct {
		*firewall.Context
		Valid bool `json:"valid"`
	}{
		Context: authContext,
		Valid:   true,
	})
}

func (h *WardenHandler) Allowed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = herodot.NewContext()
	if _, err := h.Warden.TokenAllowed(ctx, h.Warden.TokenFromRequest(r), &ladon.Request{
		Resource: "rn:hydra:warden:allowed",
		Action:   "decide",
	}, "hydra.warden"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	var access = new(ladon.Request)
	if err := json.NewDecoder(r.Body).Decode(access); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}
	defer r.Body.Close()

	if err := h.Warden.IsAllowed(ctx, access); err != nil {
		h.H.Write(ctx, w, r, &notAllowed)
		return
	}

	res := notAllowed
	res.Allowed = true
	h.H.Write(ctx, w, r, &res)
}

func (h *WardenHandler) TokenAllowed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := herodot.NewContext()
	clientCtx, err := h.Warden.TokenAllowed(ctx, h.Warden.TokenFromRequest(r), &ladon.Request{
		Resource: "rn:hydra:warden:token:allowed",
		Action:   "decide",
	}, "warden.token")
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	var ar = WardenAccessRequest{
		Request:                 new(ladon.Request),
		WardenAuthorizedRequest: new(WardenAuthorizedRequest),
	}
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}
	defer r.Body.Close()

	authContext, err := h.Warden.TokenAllowed(ctx, ar.Token, ar.Request, ar.Scopes...)
	if err != nil {
		h.H.Write(ctx, w, r, &notAllowed)
		return
	}

	authContext.Audience = clientCtx.Subject
	h.H.Write(ctx, w, r, struct {
		*firewall.Context
		Allowed bool `json:"allowed"`
	}{
		Context: authContext,
		Allowed: true,
	})
}

func TokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		return ""
	}

	return split[1]
}
