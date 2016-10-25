package jwk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"golang.org/x/net/context"
)

type Handler struct {
	Manager    Manager
	Generators map[string]KeyGenerator
	H          herodot.Herodot
	W          firewall.Firewall
}

func (h *Handler) GetGenerators() map[string]KeyGenerator {
	if h.Generators == nil || len(h.Generators) == 0 {
		h.Generators = map[string]KeyGenerator{
			"RS256": &RS256Generator{},
			"ES521": &ECDSA521Generator{},
			"HS256": &HS256Generator{
				Length: 32,
			},
		}
	}
	return h.Generators
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST("/keys/:set", h.Create)
	r.PUT("/keys/:set", h.UpdateKeySet)
	r.GET("/keys/:set", h.GetKeySet)
	r.DELETE("/keys/:set", h.DeleteKeySet)

	r.PUT("/keys/:set/:key", h.UpdateKey)
	r.GET("/keys/:set/:key", h.GetKey)
	r.DELETE("/keys/:set/:key", h.DeleteKey)

}

type createRequest struct {
	Algorithm string `json:"alg"`
	KeyID     string `json:"id"`
}

type joseWebKeySetRequest struct {
	Keys []json.RawMessage `json:"keys"`
}

func (h *Handler) DeleteKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var setName = ps.ByName("set")
	var keyName = ps.ByName("key")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:keys:" + setName + ":" + keyName,
		Action:   "delete",
	}, "hydra.keys.delete"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.DeleteKey(setName, keyName); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var setName = ps.ByName("set")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:keys:" + setName,
		Action:   "delete",
	}, "hydra.keys.delete"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.DeleteKeySet(setName); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var keyRequest createRequest
	var set = ps.ByName("set")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:keys:" + set,
		Action:   "create",
	}, "hydra.keys.create"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&keyRequest); err != nil {
		h.H.WriteError(ctx, w, r, errors.Wrap(err, ""))
	}

	generator, found := h.GetGenerators()[keyRequest.Algorithm]
	if !found {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.Errorf("Generator %s unknown", keyRequest.Algorithm))
		return
	}

	keys, err := generator.Generate(keyRequest.KeyID)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.AddKeySet(set, keys); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.WriteCreated(ctx, w, r, fmt.Sprintf("%s://%s/keys/%s", r.URL.Scheme, r.URL.Host, set), keys)
}

func (h *Handler) UpdateKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var requests joseWebKeySetRequest
	var keySet = new(jose.JsonWebKeySet)
	var set = ps.ByName("set")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:keys:" + set,
		Action:   "update",
	}, "hydra.keys.update"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		h.H.WriteError(ctx, w, r, errors.Wrap(err, ""))
		return
	}

	for _, request := range requests.Keys {
		key := &jose.JsonWebKey{}
		if err := key.UnmarshalJSON(request); err != nil {
			h.H.WriteError(ctx, w, r, errors.Wrap(err, ""))
		}
		keySet.Keys = append(keySet.Keys, *key)
	}

	if err := h.Manager.AddKeySet(set, keySet); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, keySet)
}

func (h *Handler) UpdateKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var key jose.JsonWebKey
	var set = ps.ByName("set")

	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		h.H.WriteError(ctx, w, r, errors.Wrap(err, ""))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:keys:" + set + ":" + key.KeyID,
		Action:   "update",
	}, "hydra.keys.update"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.AddKey(set, &key); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, key)
}

func (h *Handler) GetKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var setName = ps.ByName("set")
	var keyName = ps.ByName("key")

	if err := h.W.IsAllowed(ctx, &firewall.AccessRequest{
		Subject:  "",
		Resource: "rn:hydra:keys:" + setName + ":" + keyName,
		Action:   "get",
	}); err == nil {
		// Allow unauthorized requests to access this resource if it is enabled by policies
	} else if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: "rn:hydra:keys:" + setName + ":" + keyName,
		Action:   "get",
	}, "hydra.keys.get"); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	keys, err := h.Manager.GetKey(setName, keyName)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, keys)
}

func (h *Handler) GetKeySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var setName = ps.ByName("set")

	keys, err := h.Manager.GetKeySet(setName)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	for _, key := range keys.Keys {
		if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
			Resource: "rn:hydra:keys:" + setName + ":" + key.KeyID,
			Action:   "get",
		}, "hydra.keys.get"); err != nil {
			h.H.WriteError(ctx, w, r, err)
			return
		}
	}

	h.H.Write(ctx, w, r, keys)
}
