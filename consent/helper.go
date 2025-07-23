// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"net/url"
	"strings"

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

func matchScopes(scopeStrategy fosite.ScopeStrategy, grantedScope, requestedScope []string) bool {
	for _, scope := range requestedScope {
		if !scopeStrategy(grantedScope, scope) {
			return false
		}
	}
	return true
}

func caseInsensitiveFilterParam(q url.Values, key string) url.Values {
	query := url.Values{}
	key = strings.ToLower(key)
	for k, vs := range q {
		if key == strings.ToLower(k) {
			query.Set(k, "****")
		} else {
			for _, v := range vs {
				query.Add(k, v)
			}
		}
	}
	return query
}
