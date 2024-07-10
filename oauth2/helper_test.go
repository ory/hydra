// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
)

func Tokens(c fosite.Configurator, length int) (res [][]string) {
	s := &oauth2.HMACSHAStrategy{BaseHMACSHAStrategy: &oauth2.BaseHMACSHAStrategy{Enigma: &hmac.HMACStrategy{Config: c}, Config: c}}

	for i := 0; i < length; i++ {
		tok, sig, _ := s.Enigma.Generate(context.Background())
		res = append(res, []string{sig, tok})
	}
	return res
}
