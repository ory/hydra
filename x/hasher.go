/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */
package x

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/x/hasherx"

	"go.opentelemetry.io/otel"

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

func (b *Hasher) Hash(ctx context.Context, data []byte) ([]byte, error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "x.hasher.Hash")
	defer span.End()

	switch b.c.GetHasherAlgorithm(ctx) {
	case HashAlgorithmBCrypt:
		return b.bcrypt.Generate(ctx, data)
	case HashAlgorithmPBKDF2:
		fallthrough
	default:
		return b.pbkdf2.Generate(ctx, data)
	}
}

func (b *Hasher) Compare(ctx context.Context, hash, data []byte) error {
	_, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "x.hasher.Hash")
	defer span.End()

	if err := hasherx.Compare(ctx, data, hash); err != nil {
		return errorsx.WithStack(err)
	}
	return nil
}
