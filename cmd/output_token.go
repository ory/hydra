package cmd

import (
	"fmt"
	"github.com/ory/x/stringsx"
	"time"

	"golang.org/x/oauth2"
)

type (
	outputOAuth2Token oauth2.Token
)

func (_ outputOAuth2Token) Header() []string {
	return []string{"ACCESS TOKEN", "REFRESH TOKEN", "ID TOKEN", "EXPIRY"}
}

func (i outputOAuth2Token) Columns() []string {
	token := oauth2.Token(i)
	return []string{
		i.AccessToken,
		stringsx.Coalesce(i.RefreshToken, "<empty>"),
		stringsx.Coalesce(fmt.Sprintf("%s", token.Extra("id_token")), "<empty>"),
		i.Expiry.Round(time.Second).String(),
	}
}

func (i outputOAuth2Token) Interface() interface{} {
	return i
}
