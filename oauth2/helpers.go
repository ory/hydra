// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	gofrsuuid "github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/x"
)

func signatureFromJTI(jti string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(jti)))
}

type BlacklistedJTI struct {
	JTI    string         `db:"-"`
	ID     string         `db:"signature"`
	Expiry time.Time      `db:"expires_at"`
	NID    gofrsuuid.UUID `db:"nid"`
}

func (j *BlacklistedJTI) AfterFind(_ *pop.Connection) error {
	j.Expiry = j.Expiry.UTC()
	return nil
}

func (BlacklistedJTI) TableName() string {
	return "hydra_oauth2_jti_blacklist"
}

func NewBlacklistedJTI(jti string, exp time.Time) *BlacklistedJTI {
	return &BlacklistedJTI{
		JTI: jti,
		ID:  signatureFromJTI(jti),
		// because the database timestamp types are not as accurate as time.Time we truncate to seconds (which should always work)
		Expiry: exp.UTC().Truncate(time.Second),
	}
}

type AssertionJWTReader interface {
	x.FositeStorer
	GetClientAssertionJWT(ctx context.Context, jti string) (*BlacklistedJTI, error)
	SetClientAssertionJWTRaw(context.Context, *BlacklistedJTI) error
}
