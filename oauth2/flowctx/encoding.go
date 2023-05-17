package flowctx

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"

	"github.com/ory/hydra/v2/jwk"
	"github.com/pkg/errors"
)

// Decode decodes the given string to a value.
func Decode[T any](ctx context.Context, cipher *jwk.AEAD, encoded string) (*T, error) {
	plaintext, err := cipher.Decrypt(ctx, encoded)
	if err != nil {
		return nil, err
	}

	rawBytes, err := gzip.NewReader(bytes.NewReader(plaintext))
	if err != nil {
		return nil, err
	}
	defer func() { _ = rawBytes.Close() }()

	var val T
	if err = json.NewDecoder(rawBytes).Decode(&val); err != nil {
		return nil, err
	}

	return &val, nil
}

// Encode encodes the given value to a string.
func Encode(ctx context.Context, cipher *jwk.AEAD, val any) (s string, err error) {
	// Steps:
	// 1. Encode to JSON
	// 2. GZIP
	// 3. Encrypt with AEAD (AES-GCM) + Base64 URL-encode
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)

	if err = json.NewEncoder(gz).Encode(val); err != nil {
		return "", err
	}
	if err = gz.Close(); err != nil {
		return "", err
	}

	return cipher.Encrypt(ctx, b.Bytes())
}

// EncodeFromContext encodes the value stored in the context under the given cookie name.
func EncodeFromContext(ctx context.Context, cipher *jwk.AEAD, cookieName string) (s string, err error) {
	v, ok := ctx.Value(contextKey(cookieName)).(*Value)
	if !ok || v == nil {
		return "", errors.WithStack(ErrNoValueInCtx)
	}

	return Encode(ctx, cipher, v.Ptr)
}

// SetCookie looks up the value stored in the context under the given cookie name and sets it as a cookie on the
// response writer.
func SetCookie(ctx context.Context, w http.ResponseWriter, cipher *jwk.AEAD, cookieName string) error {
	v, ok := ctx.Value(contextKey(cookieName)).(*Value)
	if !ok || v == nil {
		return errors.WithStack(ErrNoValueInCtx)
	}

	cookie, err := Encode(ctx, cipher, v.Ptr)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    cookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
