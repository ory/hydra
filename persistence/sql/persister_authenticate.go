// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
)

func (p *Persister) Authenticate(ctx context.Context, name, secret string) (subject string, err error) {
	session, err := p.r.Kratos().Authenticate(ctx, name, secret)
	if err != nil {
		return "", err
	}
	return session.Identity.Id, nil
}
