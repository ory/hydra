// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
)

func sanitizeClientFromRequest(ar fosite.AuthorizeRequester) *client.Client {
	return sanitizeClient(ar.GetClient().(*client.Client))
}

func sanitizeClient(c *client.Client) *client.Client {
	cc := new(client.Client)
	// Remove the hashed secret here
	*cc = *c
	cc.Secret = ""
	return cc
}

func matchScopes(scopeStrategy fosite.ScopeStrategy, previousConsent []AcceptOAuth2ConsentRequest, requestedScope []string) *AcceptOAuth2ConsentRequest {
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
