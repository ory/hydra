package warden

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
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
	Warden Warden
	Ladon  ladon.Warden
}

type WardenResponse struct {
	*Context
}

type WardenAuthorizedRequest struct {
	Scopes       []string `json:"scopes"`
	InspectToken string   `json:"inspectToken"`
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

	authContext, err := h.Warden.Authorized(ctx, ar.InspectToken, ar.Scopes...)
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

	authContext, err := h.Warden.ActionAllowed(ctx, ar.InspectToken, ar.Request, ar.Scopes...)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	authContext.Audience = clientCtx.Subject
	h.H.Write(ctx, w, r, authContext)
}

func (h *WardenHandler) authorizeClient(ctx context.Context, w http.ResponseWriter, r *http.Request, action string) (*Context, error) {
	authctx, err := h.Warden.ActionAllowed(ctx, TokenFromRequest(r), &ladon.Request{
		Action: action,
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
