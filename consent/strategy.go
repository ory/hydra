// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/flow"
)

var _ Strategy = new(DefaultStrategy)

type Strategy interface {
	HandleOAuth2AuthorizationRequest(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		req fosite.AuthorizeRequester,
	) (*flow.AcceptOAuth2ConsentRequest, *flow.Flow, error)
	HandleOAuth2DeviceAuthorizationRequest(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
	) (*flow.AcceptOAuth2ConsentRequest, *flow.Flow, error)
	HandleOpenIDConnectLogout(ctx context.Context, w http.ResponseWriter, r *http.Request) (*flow.LogoutResult, error)
	HandleHeadlessLogout(ctx context.Context, w http.ResponseWriter, r *http.Request, sid string) error
	ObfuscateSubjectIdentifier(ctx context.Context, cl fosite.Client, subject, forcedIdentifier string) (string, error)
}
