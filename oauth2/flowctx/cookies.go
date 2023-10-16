// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flowctx

import "github.com/ory/hydra/v2/client"

type (
	CookieSuffixer interface {
		CookieSuffix() string
	}

	StaticSuffix string
	clientID     string
)

func (s StaticSuffix) CookieSuffix() string { return string(s) }
func (s clientID) GetID() string            { return string(s) }

const (
	loginSessionCookie = "ory_hydra_loginsession"
)

func LoginSessionCookie(suffix CookieSuffixer) string {
	return loginSessionCookie + "_" + suffix.CookieSuffix()
}

func SuffixForClient(c client.IDer) StaticSuffix {
	return StaticSuffix(client.CookieSuffix(c))
}

func SuffixFromStatic(id string) StaticSuffix {
	return SuffixForClient(clientID(id))
}
