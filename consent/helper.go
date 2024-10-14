// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/flow"
)

func sanitizeClientFromRequest(ar fosite.AuthorizeRequester) *client.Client {
	return client.GetSanitizedCopy(ar.GetClient().(*client.Client))
}

func matchScopes(scopeStrategy fosite.ScopeStrategy, previousConsent []flow.AcceptOAuth2ConsentRequest, requestedScope []string) *flow.AcceptOAuth2ConsentRequest {
	for _, cs := range previousConsent {
		var found = true
		for _, scope := range requestedScope {
			if !scopeStrategy(cs.GrantedScope, scope) {
				found = false
				break
			}
		}

		if found {
			return &cs
		}
	}

	return nil
}
