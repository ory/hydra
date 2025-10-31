// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/x/hasherx"
	"github.com/ory/x/otelx"
)

var _ fosite.Hasher = (*Hasher)(nil)

// Hasher implements fosite.Hasher.
type Hasher struct {
	t      TracingProvider
	c      config
	bcrypt *hasherx.Bcrypt
	pbkdf2 *hasherx.PBKDF2
}

type config interface {
	hasherx.PBKDF2Configurator
	hasherx.BCryptConfigurator
	GetHasherAlgorithm(ctx context.Context) string
}

// NewHasher returns a new BCrypt instance.
func NewHasher(t TracingProvider, c config) *Hasher {
	return &Hasher{
		t:      t,
		c:      c,
		bcrypt: hasherx.NewHasherBcrypt(c),
		pbkdf2: hasherx.NewHasherPBKDF2(c),
	}
}

const (
	hashAlgorithmBCrypt = "bcrypt"
	hashAlgorithmPBKDF2 = "pbkdf2"
)

func (h *Hasher) Hash(ctx context.Context, data []byte) (_ []byte, err error) {
	ctx, span := h.t.Tracer(ctx).Tracer().Start(ctx, "x.hasher.Hash")
	defer otelx.End(span, &err)

	alg := h.c.GetHasherAlgorithm(ctx)
	span.SetAttributes(attribute.String("algorithm", alg))

	switch alg {
	case hashAlgorithmBCrypt:
		return h.bcrypt.Generate(ctx, data)
	case hashAlgorithmPBKDF2:
		fallthrough
	default:
		return h.pbkdf2.Generate(ctx, data)
	}
}

func (h *Hasher) Compare(ctx context.Context, hash, data []byte) (err error) {
	ctx, span := h.t.Tracer(ctx).Tracer().Start(ctx, "x.hasher.Compare")
	defer otelx.End(span, &err)

	return hasherx.Compare(ctx, data, hash)
}
