// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"

	"github.com/ory/hydra/v2/oauth2/flowctx"
)

func FromCtx(ctx context.Context) (*Flow, error) {
	return flowctx.FromCtx[Flow](ctx, flowctx.FlowCookie)
}
func SetInCtx(ctx context.Context, f *Flow) error {
	return flowctx.SetCtx(ctx, flowctx.FlowCookie, f)
}
