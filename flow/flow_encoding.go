// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/otelx"
)

type decodeDependencies interface {
	CipherProvider
	x.NetworkProvider
	config.Provider
	x.TracingProvider
}

func DecodeFromLoginChallenge(ctx context.Context, d decodeDependencies, loginChallenge string) (_ *Flow, err error) {
	ctx, span := d.Tracer(ctx).Tracer().Start(ctx, "flow.DecodeFromLoginChallenge")
	defer otelx.End(span, &err)

	f, err := Decode[Flow](ctx, d.FlowCipher(), loginChallenge, AsLoginChallenge)
	if err != nil {
		return nil, errors.WithStack(x.ErrNotFound.WithWrap(err))
	}

	if f.NID != d.Networker().NetworkID(ctx) {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	if f.RequestedAt.Add(d.Config().ConsentRequestMaxAge(ctx)).Before(time.Now()) {
		return nil, errors.WithStack(fosite.ErrRequestUnauthorized.WithHint("The login request has expired, please try again."))
	}

	return f, nil
}
