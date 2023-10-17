// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import "context"

func (p *Persister) Authenticate(ctx context.Context, name, secret string) error {
	return p.r.Kratos().Authenticate(ctx, name, secret)
}
