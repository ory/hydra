package cmd

import (
	"fmt"
	"time"

	hydra "github.com/ory/hydra-client-go"
)

type (
	outputOAuth2TokenIntrospection hydra.IntrospectedOAuth2Token
)

func (_ outputOAuth2TokenIntrospection) Header() []string {
	return []string{"ACTIVE", "SUBJECT", "CLIENT ID", "SCOPE", "EXPIRY", "TOKEN USE"}
}

func (i outputOAuth2TokenIntrospection) Columns() []string {
	return []string{
		fmt.Sprintf("%v", i.Active),
		fmt.Sprintf("%v", i.Sub),
		fmt.Sprintf("%v", i.ClientId),
		fmt.Sprintf("%v", i.Scope),
		fmt.Sprintf("%v", i.Scope),
		fmt.Sprintf("%v", i.TokenUse),
		fmt.Sprintf("%v", time.Unix(*i.Exp, 0).String()),
	}
}

func (i outputOAuth2TokenIntrospection) Interface() interface{} {
	return i
}
