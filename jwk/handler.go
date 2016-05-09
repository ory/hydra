package jwk

import (
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/go-errors/errors"
	"fmt"
)

type Handler struct {
	Manager Manager
	Generators map[string]KeyGenerator
	H       herodot.Herodot
	W       warden.Warden
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/keys/:set", h.Create)
	r.GET("/keys/:set", h.GetSet)
	r.GET("/keys/:set/:id", h.GetKey)
}

type createRequest struct {
	Algorithm string `json:"alg"`
	KeyID string `json:"id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var keyRequest createRequest
	var ctx = context.Background()
	if err := json.NewDecoder(r.Body).Decode(&keyRequest); err != nil {
		h.H.WriteError(ctx, w, r, err)
	}

	generator, found := h.Generators[keyRequest.Algorithm]
	if !found {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.Errorf("Generator %s unknown", keyRequest.Algorithm))
		return
	}

	keys, err := generator.Generate(keyRequest.KeyID)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	set := ps.ByName("set")
	if err := h.Manager.AddKeys(set); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.WriteCreated(ctx, w, r, fmt.Sprintf("%s://%s/keys/%s", r.URL.Scheme, r.URL.Host, set), keys)
}
