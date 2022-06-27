package cmd

import (
	"fmt"

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
		i.RefreshToken,
		fmt.Sprintf("%v", token.Extra("id_token")),
		i.Expiry.String(),
	}
}

func (i outputOAuth2Token) Interface() interface{} {
	return i
}
