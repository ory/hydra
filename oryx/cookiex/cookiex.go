// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package cookiex encodes authenticated, encrypted HTTP cookies. Values are
// sealed with the Ory-wide AEAD (see ory/x/aead) under keys derived from the
// configured secrets with HKDF-SHA256, and bound to their purpose and cookie
// name. The package replaces the gorilla/securecookie-based cookie stores; a
// decode-only fallback reads cookies minted by those stores until they have
// expired.
package cookiex

import (
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/aead"
)

const (
	// formatPrefix marks the versioned wire format. A dot is not part of the
	// base64url alphabet, so legacy securecookie values can never collide
	// with the prefix.
	formatPrefix = "v1."
	// kdfInfo domain-separates the HKDF key derivation.
	kdfInfo = "ory/x/cookiex/v1"
	// aadPrefix domain-separates the additional authenticated data.
	aadPrefix = "ory/x/cookiex/v1"
	// defaultMaxAge matches the gorilla/securecookie default that the
	// previous cookie stores relied on.
	defaultMaxAge = 30 * 24 * time.Hour
	// maxCookieValueLength is the browser limit that securecookie also
	// enforced.
	maxCookieValueLength = 4096
)

// ErrInvalidCookie is returned when a cookie with the requested name is
// present but none of its values could be decoded (or none matched).
var ErrInvalidCookie = errors.New("cookiex: cookie could not be decoded")

type (
	// Option configures a Codec.
	Option func(*config)

	config struct {
		maxAge         time.Duration
		legacyKeyPairs [][]byte
		legacyEncode   bool
	}
)

// WithMaxAge overrides how old a cookie may be before decoding rejects it.
// The default is 30 days; zero disables the check. The value must not be
// negative.
func WithMaxAge(d time.Duration) Option {
	return func(c *config) { c.maxAge = d }
}

// Codec seals and opens cookie values of type T. T is serialized as JSON.
// A Codec is safe for concurrent use.
type Codec[T any] struct {
	purpose string
	keys    [][32]byte
	maxAge  time.Duration
	legacy  legacyState
	now     func() time.Time
}

// New returns a codec for the given purpose. The purpose is bound into the
// ciphertext and used as the metric label; it must be a short constant like
// "kratos/session". Because the purpose is embedded in the additional
// authenticated data, it must be non-empty and must not contain a pipe
// character. The codec seals with a key derived from the first secret
// and opens with keys derived from any of them, so secrets rotate by
// prepending a new one.
func New[T any](purpose string, secrets [][]byte, opts ...Option) (*Codec[T], error) {
	if purpose == "" || strings.Contains(purpose, "|") {
		return nil, errors.New("cookiex: purpose must be non-empty and must not contain a pipe character")
	}
	if len(secrets) == 0 {
		return nil, errors.New("cookiex: at least one secret is required")
	}
	cfg := config{maxAge: defaultMaxAge}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.maxAge < 0 {
		return nil, errors.New("cookiex: max age must not be negative")
	}
	if cfg.legacyEncode && len(cfg.legacyKeyPairs) == 0 {
		return nil, errors.New("cookiex: legacy encode requires legacy key pairs")
	}
	keys := make([][32]byte, len(secrets))
	for i, secret := range secrets {
		key, err := hkdf.Key(sha256.New, secret, nil, kdfInfo, 32)
		if err != nil {
			return nil, errors.Wrap(err, "cookiex: cannot derive key")
		}
		keys[i] = [32]byte(key)
	}
	return &Codec[T]{
		purpose: purpose,
		keys:    keys,
		maxAge:  cfg.maxAge,
		legacy:  newLegacyState(cfg, cfg.maxAge),
		now:     time.Now,
	}, nil
}

// envelope wraps the JSON payload with the seal time so decoding can enforce
// the max age.
type envelope struct {
	IssuedAt int64           `json:"iat"`
	Values   json.RawMessage `json:"v"`
}

// aad binds a ciphertext to this codec's purpose and the cookie name, so a
// sealed value cannot be replayed as a different cookie or in a different
// context, even under the same key.
func (c *Codec[T]) aad(name string) []byte {
	return []byte(aadPrefix + "|" + c.purpose + "|" + name)
}

func (c *Codec[T]) seal(name string, value T) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "cookiex: cannot marshal cookie value")
	}
	plaintext, err := json.Marshal(envelope{IssuedAt: c.now().Unix(), Values: payload})
	if err != nil {
		return "", errors.Wrap(err, "cookiex: cannot marshal envelope")
	}
	a, err := aead.New(c.keys[0])
	if err != nil {
		return "", errors.Wrap(err, "cookiex: cannot create AEAD")
	}
	// The nonce is prepended to the ciphertext. AEADs that manage the nonce
	// internally report a nonce size of zero, so this also covers them.
	nonce := make([]byte, a.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", errors.Wrap(err, "cookiex: cannot generate nonce")
	}
	sealed := a.Seal(nonce, nonce, plaintext, c.aad(name))
	return formatPrefix + base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (c *Codec[T]) open(name, value string) (T, error) {
	var zero T
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(value, formatPrefix))
	if err != nil {
		return zero, errors.WithStack(ErrInvalidCookie)
	}
	for _, key := range c.keys {
		a, err := aead.New(key)
		if err != nil {
			return zero, errors.Wrap(err, "cookiex: cannot create AEAD")
		}
		if len(raw) < a.NonceSize() {
			continue
		}
		plaintext, err := a.Open(nil, raw[:a.NonceSize()], raw[a.NonceSize():], c.aad(name))
		if err != nil {
			continue
		}
		var env envelope
		if err := json.Unmarshal(plaintext, &env); err != nil {
			return zero, errors.WithStack(ErrInvalidCookie)
		}
		if c.maxAge > 0 && c.now().Sub(time.Unix(env.IssuedAt, 0)) > c.maxAge {
			return zero, errors.WithStack(ErrInvalidCookie)
		}
		var out T
		if err := json.Unmarshal(env.Values, &out); err != nil {
			return zero, errors.WithStack(ErrInvalidCookie)
		}
		return out, nil
	}
	return zero, errors.WithStack(ErrInvalidCookie)
}

// Set seals value into cookie.Value and writes the cookie to w. All other
// attributes (name, path, domain, max-age, secure, http-only, same-site) must
// be set by the caller on the cookie.
func (c *Codec[T]) Set(w http.ResponseWriter, cookie *http.Cookie, value T) error {
	var encoded string
	var err error
	if c.legacy.encode {
		encoded, err = c.sealLegacy(cookie.Name, value)
	} else {
		encoded, err = c.seal(cookie.Name, value)
	}
	if err != nil {
		return err
	}
	if len(encoded) > maxCookieValueLength {
		return errors.Errorf("cookiex: encoded cookie %q exceeds %d bytes", cookie.Name, maxCookieValueLength)
	}
	cookie.Value = encoded
	http.SetCookie(w, cookie)
	return nil
}

// Get returns the value of the first cookie named name that decodes. It
// returns http.ErrNoCookie when no cookie with that name is present, and
// ErrInvalidCookie when none of the present values decode.
func (c *Codec[T]) Get(r *http.Request, name string) (T, error) {
	return c.GetMatching(r, name, nil)
}

// GetMatching returns the value of the first cookie named name that decodes
// and matches. Requests can carry multiple cookies with the same name, for
// example one from a parent and one from a sub domain; match selects among
// them. A nil match matches any value. It returns http.ErrNoCookie when no
// cookie with that name is present, and ErrInvalidCookie when none of the
// present values decode and match.
func (c *Codec[T]) GetMatching(r *http.Request, name string, match func(T) bool) (T, error) {
	if match == nil {
		match = func(T) bool { return true }
	}
	var zero T
	found := false
	for _, ck := range r.Cookies() {
		if ck.Name != name {
			continue
		}
		found = true
		// Refuse oversized values before doing any decode work, matching the
		// limit Set enforces and securecookie's own decode-side check.
		if len(ck.Value) > maxCookieValueLength {
			continue
		}
		var value T
		var err error
		if strings.HasPrefix(ck.Value, formatPrefix) {
			value, err = c.open(name, ck.Value)
		} else {
			value, err = c.openLegacy(name, ck.Value)
		}
		if err != nil {
			continue
		}
		if match(value) {
			return value, nil
		}
	}
	if !found {
		return zero, errors.WithStack(http.ErrNoCookie)
	}
	return zero, errors.WithStack(ErrInvalidCookie)
}

// Expire overwrites the cookie with an expired empty one, deleting it from
// the client.
func (c *Codec[T]) Expire(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(1, 0)
	http.SetCookie(w, cookie)
}
