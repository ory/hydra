package herodot

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

type jsonError struct {
	RequestID string `json:"request"`
	Message   string `json:"message"`
	*Error
}

type JSON struct {
	Logger logrus.FieldLogger
}

func (h *JSON) WriteCreated(ctx context.Context, w http.ResponseWriter, r *http.Request, location string, e interface{}) {
	w.Header().Set("Location", location)
	h.WriteCode(ctx, w, r, http.StatusCreated, e)
}

func (h *JSON) Write(ctx context.Context, w http.ResponseWriter, r *http.Request, e interface{}) {
	h.WriteCode(ctx, w, r, http.StatusOK, e)
}

func (h *JSON) WriteCode(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, e interface{}) {
	js, err := json.Marshal(e)
	if err != nil {
		h.WriteError(ctx, w, r, err)
		return
	}

	if code == 0 {
		code = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func (h *JSON) WriteError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	e := ToError(err)
	h.WriteErrorCode(ctx, w, r, e.StatusCode, err)
	return
}

func (h *JSON) WriteErrorCode(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, err error) {
	id, _ := ctx.Value(RequestIDKey).(string)
	if id == "" {
		id = uuid.New()
	}
	if code == 0 {
		code = http.StatusInternalServerError
	}

	LogError(err, id, code)
	je := ToError(err)
	je.StatusCode = code
	h.WriteCode(ctx, w, r, je.StatusCode, &jsonError{
		RequestID: id,
		Error:     ToError(err),
		Message:   err.Error(),
	})
}
