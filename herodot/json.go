package herodot

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

type jsonError struct {
	RequestID string `json:"requestId"`
	Error     string `json:"error"`
	Code      int    `json:"code"`
}

type JSON struct {
	Logger logrus.FieldLogger
}

func (h *JSON) WriteCreated(ctx context.Context, w http.ResponseWriter, r *http.Request, e interface{}, location string) error {
	w.Header().Set("Location", location)
	return h.WriteCode(ctx, w, r, e, http.StatusCreated)
}

func (h *JSON) Write(ctx context.Context, w http.ResponseWriter, r *http.Request, e interface{}) error {
	return h.WriteCode(ctx, w, r, e, http.StatusOK)
}

func (h *JSON) WriteCode(ctx context.Context, w http.ResponseWriter, r *http.Request, e interface{}, code int) error {
	js, err := json.Marshal(e)
	if err != nil {
		h.WriteError(ctx, w, r, err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
	return nil
}

func (h *JSON) WriteError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error {
	if e, ok := err.(*Error); ok {
		return h.WriteErrorCode(ctx, w, r, e, e.Code)
	}

	return h.WriteErrorCode(ctx, w, r, err, http.StatusInternalServerError)
}

func (h *JSON) WriteErrorCode(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, status int) error {
	id, _ := ctx.Value(RequestIDKey).(string)
	if id == "" {
		id = uuid.New()
	}

	if h.Logger == nil {
		h.Logger = logrus.New()
	}

	if e, ok := err.(*Error); ok {
		h.Logger.WithError(e.Error).WithField("request_id", id).WithField("status", status).WithField("stack", e.ErrorStack())
	} else {
		h.Logger.WithError(e.Error).WithField("request_id", id).WithField("status", status)
	}

	return h.Write(ctx, w, r, &jsonError{
		RequestID: id,
		Error:     err.Error(),
		Code:      status,
	})
}
