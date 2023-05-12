package flowctx

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/errorsx"
	"github.com/pkg/errors"
)

type (
	contextKey string

	Value struct {
		Ptr     any
		ptrType reflect.Type
		decoded []byte
	}

	Dependencies interface {
		KeyCipher() *jwk.AEAD
		x.RegistryWriter
		x.RegistryLogger
	}

	Middleware struct {
		cookieName     string
		postDecodeHook func(ctx context.Context, val any) error
		Dependencies
	}

	fromCtxOptions struct {
		postDecodeHook func(ctx context.Context, val any) error
	}

	FromCtxOpt func(o *fromCtxOptions)
)

var (
	ErrCookieCorrupted = fosite.ErrInvalidRequest.WithHint("cookie corrupted")
	ErrNoValueInCtx    = errors.New("no value in context")
)

func WithPostDecodeHook(hook func(context.Context, any) error) FromCtxOpt {
	return func(o *fromCtxOptions) {
		o.postDecodeHook = hook
	}
}

func NewMiddleware(cookieName string, dependencies Dependencies) *Middleware {
	m := &Middleware{
		cookieName:   cookieName,
		Dependencies: dependencies,
	}

	return m
}

// FromCtx returns the underlying value from the context. If the value is nil, the second return value is false.
func FromCtx[T any](ctx context.Context, cookieName string, opt ...FromCtxOpt) (*T, error) {
	opts := newFromCtxOpts(opt...)

	v, ok := ctx.Value(contextKey(cookieName)).(*Value)
	if !ok || v == nil {
		return nil, errors.WithStack(ErrNoValueInCtx)
	}

	var t T
	prtType := reflect.TypeOf(t)
	if v.ptrType != nil && v.ptrType != prtType {
		return nil, errors.Errorf("expected type %q but got %q", v.ptrType.String(), prtType.String())
	}
	v.ptrType = prtType

	switch {
	case v.Ptr == nil && v.decoded != nil:
		// Value was decoded before, but not yet unmarshaled.
		if err := json.Unmarshal(v.decoded, &t); err != nil {
			return nil, err
		}
		if opts.postDecodeHook != nil {
			if err := opts.postDecodeHook(ctx, &t); err != nil {
				return nil, err
			}
		}
		v.Ptr = &t

		return &t, nil

	case v.Ptr == nil:
		// Value was never set from cookie before.
		v.Ptr = &t

		return &t, nil

	default:
		return v.Ptr.(*T), nil
	}
}

func SetCtx[T any](ctx context.Context, cookieName string, val *T) error {
	v, ok := ctx.Value(contextKey(cookieName)).(*Value)
	if !ok || v == nil {
		return errors.WithStack(ErrNoValueInCtx)
	}

	prtType := reflect.TypeOf(*val)
	if v.ptrType != nil && v.ptrType != prtType {
		return errors.Errorf("expected type %q but got %q", v.ptrType.String(), prtType.String())
	}
	v.ptrType = prtType
	v.Ptr = val

	return nil
}

func (m *Middleware) Handle(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx, value, err := m.fromHTTP(r.Context(), r)
		if err != nil {
			m.Dependencies.Writer().WriteError(w, r, errorsx.WithStack(ErrCookieCorrupted))
			return
		}

		next(w, r.WithContext(ctx), params)

		if err := m.setCookie(ctx, value, w); err != nil {
			m.Dependencies.Logger().WithError(err).Errorf("could not write %q cookie", m.cookieName)
		}
	}
}

func (m *Middleware) setCookie(ctx context.Context, v *Value, w http.ResponseWriter) error {
	if v.Ptr == nil {
		return nil
	}
	cookie, err := m.encode(ctx, v.Ptr)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  m.cookieName,
		Value: cookie,
	})

	return nil
}

func (m *Middleware) fromHTTP(ctx context.Context, r *http.Request) (context.Context, *Value, error) {
	var (
		v = &Value{}
	)
	ctx = context.WithValue(ctx, contextKey(m.cookieName), v)

	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return ctx, v, nil // missing cookie is not an error
	}

	v.decoded, err = m.decode(ctx, cookie.Value)
	if err != nil {
		return ctx, v, err // corrupted cookie is an error
	}

	return ctx, v, nil
}

func (m *Middleware) decode(ctx context.Context, cookie string) ([]byte, error) {
	plaintext, err := m.Dependencies.KeyCipher().Decrypt(ctx, cookie)
	if err != nil {
		return nil, err
	}

	rawBytes, err := gzip.NewReader(bytes.NewReader(plaintext))
	if err != nil {
		return nil, err
	}
	defer rawBytes.Close()

	return io.ReadAll(rawBytes)
}

func (m *Middleware) encode(ctx context.Context, t any) (s string, err error) {
	// Steps:
	// 1. Encode to JSON
	// 2. GZIP
	// 3. Encrypt with AEAD (AES-GCM) + Base64 URL-encode
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)

	if err = json.NewEncoder(gz).Encode(t); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}

	return m.Dependencies.KeyCipher().Encrypt(ctx, b.Bytes())
}

func newFromCtxOpts(opt ...FromCtxOpt) *fromCtxOptions {
	o := &fromCtxOptions{}
	for _, f := range opt {
		f(o)
	}
	return o
}
