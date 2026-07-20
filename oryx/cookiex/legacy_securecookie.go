// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cookiex

// This file isolates the gorilla/securecookie dependency. It exists so that
// cookies sealed by the previous securecookie-based stores keep decoding
// during the migration, and (behind WithLegacyEncode) so that freshly minted
// cookies stay readable by old pods during the stage-1 rollout. Delete it,
// and the securecookie dependency, once ory_legacy_cookie_decodes_total has
// flatlined fleet-wide for a full max-cookie-lifespan window.

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// legacyDecodes counts cookies that were successfully decoded via the legacy
// securecookie fallback. It is registered on the Prometheus default
// registerer, which Kratos and Hydra expose on /metrics/prometheus.
var legacyDecodes = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ory_legacy_cookie_decodes_total",
	Help: "Number of cookies successfully decoded via the legacy securecookie fallback, by codec purpose.",
}, []string{"purpose"})

// LegacySecureCookieKeyPairs derives the (hashKey, blockKey) pairs the
// securecookie-based cookie stores used: each secret signs, and its SHA-256
// digest encrypts.
func LegacySecureCookieKeyPairs(secrets [][]byte) [][]byte {
	pairs := make([][]byte, 0, 2*len(secrets))
	for _, secret := range secrets {
		block := sha256.Sum256(secret)
		pairs = append(pairs, secret, block[:])
	}
	return pairs
}

// WithLegacyKeyPairs enables decoding cookies in the legacy securecookie
// format, sealed under the given (hashKey, blockKey) pairs. While this option
// is in use, T must be unmarshalable from a flat JSON object with string
// values; WithLegacyEncode additionally requires T to marshal to one. Remove
// once all cookies sealed in the legacy format have expired.
func WithLegacyKeyPairs(pairs ...[]byte) Option {
	return func(c *config) { c.legacyKeyPairs = pairs }
}

type legacyState struct {
	codecs []securecookie.Codec
	encode bool
}

func newLegacyState(cfg config, maxAge time.Duration) legacyState {
	codecs := securecookie.CodecsFromPairs(cfg.legacyKeyPairs...)
	for _, codec := range codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(int(maxAge / time.Second))
		}
	}
	return legacyState{codecs: codecs, encode: cfg.legacyEncode}
}

// openLegacy decodes a securecookie-sealed value and bridges it into T via
// its JSON representation: legacy payloads are flat string-to-string maps, so
// string, time.Time, and UUID fields all round-trip through their JSON string
// forms.
func (c *Codec[T]) openLegacy(name, value string) (T, error) {
	var zero T
	if len(c.legacy.codecs) == 0 {
		return zero, errors.WithStack(ErrInvalidCookie)
	}
	values := make(map[any]any)
	if err := securecookie.DecodeMulti(name, value, &values, c.legacy.codecs...); err != nil {
		return zero, errors.WithStack(ErrInvalidCookie)
	}
	flat := make(map[string]string, len(values))
	for k, v := range values {
		key, keyOK := k.(string)
		val, valOK := v.(string)
		if !keyOK || !valOK {
			return zero, errors.WithStack(ErrInvalidCookie)
		}
		flat[key] = val
	}
	bridge, err := json.Marshal(flat)
	if err != nil {
		return zero, errors.Wrap(err, "cookiex: cannot bridge legacy cookie")
	}
	var out T
	if err := json.Unmarshal(bridge, &out); err != nil {
		return zero, errors.WithStack(ErrInvalidCookie)
	}
	legacyDecodes.WithLabelValues(c.purpose).Inc()
	return out, nil
}

// WithLegacyEncode makes Set seal in the legacy securecookie format under the
// first legacy key pair, so pods that only understand the legacy format can
// read freshly minted cookies during a rolling deploy. Requires
// WithLegacyKeyPairs. This is stage 1 of the rollout; a follow-up removes the
// option, flipping encoding to the v1 format. JSON null values inside the
// payload are coerced to empty strings by the bridge; do not use pointer-typed
// fields while legacy encode is enabled.
func WithLegacyEncode() Option {
	return func(c *config) { c.legacyEncode = true }
}

// sealLegacy bridges T through its JSON representation into the flat
// string-to-string map that the securecookie stores used.
func (c *Codec[T]) sealLegacy(name string, value T) (string, error) {
	buf, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "cookiex: cannot marshal cookie value")
	}
	var flat map[string]string
	if err := json.Unmarshal(buf, &flat); err != nil {
		return "", errors.Wrap(err, "cookiex: payload must be a flat JSON object with string values while legacy encode is enabled")
	}
	if flat == nil {
		return "", errors.New("cookiex: payload must be a flat JSON object with string values while legacy encode is enabled")
	}
	values := make(map[any]any, len(flat))
	for k, v := range flat {
		values[k] = v
	}
	encoded, err := securecookie.EncodeMulti(name, values, c.legacy.codecs[0])
	if err != nil {
		return "", errors.Wrap(err, "cookiex: cannot encode legacy cookie")
	}
	return encoded, nil
}
