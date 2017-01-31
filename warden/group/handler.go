package group

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/pkg/errors"
)

type membersRequest struct {
	Members []string `json:"members"`
}

type Handler struct {
	Manager Manager
	H       herodot.Herodot
	W       firewall.Firewall
}

const (
	GroupsHandlerPath = "/warden/groups"
)

const (
	GroupsResource = "rn:hydra:warden:clients"
	GroupResource  = "rn:hydra:warden:clients:%s"
	Scope          = "hydra.warden.groups"
)

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(GroupsHandlerPath, h.CreateGroup)
	r.GET(GroupsHandlerPath, h.FindGroupNames)
	r.GET(GroupsHandlerPath+"/:id", h.GetGroup)
	r.DELETE(GroupsHandlerPath+"/:id", h.DeleteGroup)
	r.POST(GroupsHandlerPath+"/:id/members", h.AddGroupMembers)
	r.DELETE(GroupsHandlerPath+"/:id/members", h.RemoveGroupMembers)
}

func (h *Handler) FindGroupNames(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = herodot.NewContext()
	var member = r.URL.Query().Get("member")

	g, err := h.Manager.FindGroupNames(member)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(GroupResource, member),
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, g)
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var g Group
	var ctx = herodot.NewContext()

	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		h.H.WriteError(ctx, w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: GroupsResource,
		Action:   "create",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.CreateGroup(&g); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.WriteCreated(ctx, w, r, GroupsHandlerPath+"/"+g.ID, &g)
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	g, err := h.Manager.GetGroup(id)
	if err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(GroupResource, id),
		Action:   "get",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	h.H.Write(ctx, w, r, g)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(GroupResource, id),
		Action:   "delete",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.DeleteGroup(id); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	var m membersRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.H.WriteError(ctx, w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(GroupResource, id),
		Action:   "add.member",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.AddGroupMembers(id, m.Members); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = ps.ByName("id")

	var m membersRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.H.WriteError(ctx, w, r, errors.WithStack(err))
		return
	}

	if _, err := h.W.TokenAllowed(ctx, h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(GroupResource, id),
		Action:   "remove.member",
	}, Scope); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	if err := h.Manager.RemoveGroupMembers(id, m.Members); err != nil {
		h.H.WriteError(ctx, w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
