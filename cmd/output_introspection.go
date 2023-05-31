// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"time"

	"github.com/ory/x/pointerx"

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
		i.Sub = pointerx.String("")
	}

	if i.ClientId == nil {
		i.ClientId = pointerx.String("")
	}

	if i.Scope == nil {
		i.Scope = pointerx.String("")
	}

	if i.TokenUse == nil {
		i.TokenUse = pointerx.String("")
	}

	if i.Exp == nil {
		i.Exp = pointerx.Int64(0)
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
