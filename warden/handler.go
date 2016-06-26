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
	"golang.org/x/net/context"
)

const (
	AuthorizedHandlerPath = "/warden/authorized"
	AllowedHandlerPath    = "/warden/allowed"
)

type WardenHandler struct {
	H      herodot.Herodot
	Warden firewall.Firewall
	Ladon  ladon.Warden
}

func NewHandler(c *config.Config, router *httprouter.Router) *WardenHandler {
	ctx := c.Context()

	h := &WardenHandler{
		H:      &herodot.JSON{},
		Warden: ctx.Warden,
		Ladon: &ladon.Ladon{
			Manager: ctx.LadonManager,
		},
	}
	h.SetRoutes(router)

	return h
}

type WardenResponse struct {
	*firewall.Context
}

type WardenAuthorizedRequest struct {
	Scopes    []string `json:"scopes"`
	Assertion string   `json:"assertion"`
}

type WardenAccessRequest struct {
	*ladon.Request
	*WardenAuthorizedRequest
}

func (h *WardenHandler) SetRoutes(r *httprouter.Router) {
	r.POST(AuthorizedHandlerPath, h.Authorized)
	r.POST(AllowedHandlerPath, h.Allowed)
}

func (h *WardenHandler) Authorized(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := herodot.NewContext()
	clientCtx, err := h.authorizeClient(ctx, w, r, "an:hydra:warden:authorized")
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

	authContext, err := h.Warden.Authorized(ctx, ar.Assertion, ar.Scopes...)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	authContext.Audience = clientCtx.Subject
	h.H.Write(ctx, w, r, authContext)

}

func (h *WardenHandler) Allowed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := herodot.NewContext()
	clientCtx, err := h.authorizeClient(ctx, w, r, "an:hydra:warden:allowed")
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	var ar WardenAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		h.H.WriteError(ctx, w, r, errors.New(err))
		return
	}

	authContext, err := h.Warden.ActionAllowed(ctx, ar.Assertion, ar.Request, ar.Scopes...)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	authContext.Audience = clientCtx.Subject
	h.H.Write(ctx, w, r, authContext)
}

func (h *WardenHandler) authorizeClient(ctx context.Context, w http.ResponseWriter, r *http.Request, action string) (*firewall.Context, error) {
	authctx, err := h.Warden.ActionAllowed(ctx, TokenFromRequest(r), &ladon.Request{
		Action:   action,
		Resource: "rn:hydra:warden",
	}, "hydra.warden")
	if err != nil {
		return nil, err
	}

	return authctx, nil
}

func TokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		return ""
	}

	return split[1]
}
