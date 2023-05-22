// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"

	"github.com/ory/hydra/v2/oauth2/flowctx"
)

func LoginSessionFromCtx(ctx context.Context) (*LoginSession, error) {
	return flowctx.FromCtx[LoginSession](ctx, flowctx.LoginSessionCookie)
}
func SetLoginSessionInCtx(ctx context.Context, l *LoginSession) error {
	return flowctx.SetCtx(ctx, flowctx.LoginSessionCookie, l)
}
