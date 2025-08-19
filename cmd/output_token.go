// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"cmp"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type (
	outputOAuth2Token oauth2.Token
)

func (outputOAuth2Token) Header() []string {
	return []string{"ACCESS TOKEN", "REFRESH TOKEN", "ID TOKEN", "EXPIRY"}
}

func (i outputOAuth2Token) Columns() []string {
	token := oauth2.Token(i)
	printIDToken := "<empty>"
	if idt := token.Extra("id_token"); idt != nil {
		printIDToken = fmt.Sprintf("%s", token.Extra("id_token"))
	}

	return []string{
		i.AccessToken,
		cmp.Or(i.RefreshToken, "<empty>"),
		printIDToken,
		i.Expiry.Round(time.Second).String(),
	}
}

func (i outputOAuth2Token) Interface() interface{} {
	return i
}
