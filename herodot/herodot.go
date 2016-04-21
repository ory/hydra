package herodot

import (
	"net/http"

	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

type key int

const RequestIDKey key = 0

func NewContext() context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, RequestIDKey, uuid.New())
}

func Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestIDKey, uuid.New())
}

type Herodot interface {
	Write(ctx context.Context, w http.ResponseWriter, r *http.Request, e interface{}) error

	WriteCode(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, e interface{}) error

	WriteCreated(ctx context.Context, w http.ResponseWriter, r *http.Request, location string, e interface{}) error

	WriteError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error

	WriteErrorCode(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, err error) error
}
