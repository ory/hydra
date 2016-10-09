package connection

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	connectionsResource = "rn:hydra:connections"
	connectionResource  = "rn:hydra:connections:%s"
	scope               = "hydra.connections"
)

type Handler struct {
	Manager Manager
	H       herodot.Herodot
	W       firewall.Firewall
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET("/connections", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.URL.Query().Get("local_subject") != "" {
			h.FindLocal(w, r, ps)
			return
		}

		if r.URL.Query().Get("remote_subject") != "" && r.URL.Query().Get("provider") != "" {
			h.FindRemote(w, r, ps)
			return
		}

		var ctx = context.Background()
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.New("Pass either [local_subject] or [remote_subject, provider] as query to this request"))
	})

	r.POST("/connections", h.Create)
	r.GET("/connections/:id", h.Get)
	r.DELETE("/connections/:id", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var conn Connection
	var ctx = context.Background()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: connectionsResource,
		Action:   "create",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, err)
		return
	}

	if v, err := govalidator.ValidateStruct(conn); err != nil {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, err)
		return
	} else if !v {
		h.H.WriteErrorCode(ctx, w, r, http.StatusBadRequest, errors.New("Payload did not validate."))
		return
	}

	conn.ID = uuid.New()
	if err := h.Manager.Create(&conn); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.WriteCreated(ctx, w, r, "/oauth2/connections/"+conn.ID, &conn)
}

func (h *Handler) FindLocal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: connectionsResource,
		Action:   "find",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	conns, err := h.Manager.FindAllByLocalSubject(r.URL.Query().Get("local_subject"))
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, conns)
}

func (h *Handler) FindRemote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: connectionsResource,
		Action:   "find",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	conns, err := h.Manager.FindByRemoteSubject(r.URL.Query().Get("provider"), r.URL.Query().Get("remote_subject"))
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, conns)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var id = ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(connectionResource, id),
		Action:   "get",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	conn, err := h.Manager.Get(id)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, conn)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = context.Background()
	var id = ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(connectionResource, id),
		Action:   "delete",
	}, scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.Delete(ps.ByName("id")); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
