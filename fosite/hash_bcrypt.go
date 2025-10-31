// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"

	"github.com/ory/x/errorsx"

	"golang.org/x/crypto/bcrypt"
)

const DefaultBCryptWorkFactor = 12

// BCrypt implements the Hasher interface by using BCrypt.
type BCrypt struct {
	Config interface {
		BCryptCostProvider
	}
}

func (b *BCrypt) Hash(ctx context.Context, data []byte) ([]byte, error) {
	wf := b.Config.GetBCryptCost(ctx)
	if wf == 0 {
		wf = DefaultBCryptWorkFactor
	}
	s, err := bcrypt.GenerateFromPassword(data, wf)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}
	return s, nil
}

func (b *BCrypt) Compare(ctx context.Context, hash, data []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, data); err != nil {
		return errorsx.WithStack(err)
	}
	return nil
}
