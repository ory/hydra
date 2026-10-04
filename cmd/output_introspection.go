// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"time"

	hydra "github.com/ory/hydra-client-go/v2"
)

type (
	outputOAuth2TokenIntrospection hydra.IntrospectedOAuth2Token
)

func (outputOAuth2TokenIntrospection) Header() []string {
	return []string{"ACTIVE", "SUBJECT", "CLIENT ID", "SCOPE", "EXPIRY", "TOKEN USE"}
}

func (i outputOAuth2TokenIntrospection) Columns() []string {
	if i.Sub == nil {
		i.Sub = new("")
	}

	if i.ClientId == nil {
		i.ClientId = new("")
	}

	if i.Scope == nil {
		i.Scope = new("")
	}

	if i.TokenUse == nil {
		i.TokenUse = new("")
	}

	if i.Exp == nil {
		i.Exp = new(int64(0))
	}

	return []string{
		fmt.Sprintf("%v", i.Active),
		fmt.Sprintf("%v", *i.Sub),
		fmt.Sprintf("%v", *i.ClientId),
		fmt.Sprintf("%v", *i.Scope),
		fmt.Sprintf("%v", time.Unix(*i.Exp, 0).String()),
		fmt.Sprintf("%v", *i.TokenUse),
	}
}

func (i outputOAuth2TokenIntrospection) Interface() interface{} {
	return i
}
