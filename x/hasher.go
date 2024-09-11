// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/x/hasherx"
	"github.com/ory/x/otelx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/x/errorsx"
)

const tracingComponent = "github.com/ory/hydra/x"

var _ fosite.Hasher = (*Hasher)(nil)

type HashAlgorithm string

func (a HashAlgorithm) String() string {
	return string(a)
}

const (
	HashAlgorithmBCrypt = HashAlgorithm("bcrypt")
	HashAlgorithmPBKDF2 = HashAlgorithm("pbkdf2")
)

// Hasher implements fosite.Hasher.
type Hasher struct {
	c      config
	bcrypt *hasherx.Bcrypt
	pbkdf2 *hasherx.PBKDF2
}

type config interface {
	hasherx.PBKDF2Configurator
	hasherx.BCryptConfigurator
	GetHasherAlgorithm(ctx context.Context) HashAlgorithm
}

// NewHasher returns a new BCrypt instance.
func NewHasher(c config) *Hasher {
	return &Hasher{
		c:      c,
		bcrypt: hasherx.NewHasherBcrypt(c),
		pbkdf2: hasherx.NewHasherPBKDF2(c),
	}
}

func (b *Hasher) Hash(ctx context.Context, data []byte) (_ []byte, err error) {
	h := b.c.GetHasherAlgorithm(ctx)

	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "x.hasher.Hash",
		trace.WithAttributes(attribute.Stringer("algorithm", h)))
	defer otelx.End(span, &err)

	switch h {
	case HashAlgorithmBCrypt:
		return b.bcrypt.Generate(ctx, data)
	case HashAlgorithmPBKDF2:
		fallthrough
	default:
		return b.pbkdf2.Generate(ctx, data)
	}
}

func (b *Hasher) Compare(ctx context.Context, hash, data []byte) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "x.hasher.Compare")
	defer otelx.End(span, &err)

	if err := hasherx.Compare(ctx, data, hash); err != nil {
		return errorsx.WithStack(err)
	}
	return nil
}
