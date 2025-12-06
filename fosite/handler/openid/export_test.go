// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"context"
	"time"

	"github.com/ory/hydra/v2/fosite"
)

func CallGenerateIDToken(ctx context.Context, lifespan time.Duration, fosr fosite.Requester, h *IDTokenHandleHelper) (token string, err error) {
	return h.generateIDToken(ctx, lifespan, fosr)
}
