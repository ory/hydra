package osin

import (
	"encoding/base64"
	"strings"

	"github.com/pborman/uuid"
)

// AuthorizeTokenGenDefault is the default authorization token generator
type AuthorizeTokenGenDefault struct {
}

func removePadding(token string) string {
	return strings.TrimRight(token, "=")
}

// GenerateAuthorizeToken generates a base64-encoded UUID code
func (a *AuthorizeTokenGenDefault) GenerateAuthorizeToken(data *AuthorizeData) (ret string, err error) {
	token := uuid.NewRandom()
	return removePadding(base64.URLEncoding.EncodeToString([]byte(token))), nil
}

// AccessTokenGenDefault is the default authorization token generator
type AccessTokenGenDefault struct {
}

// GenerateAccessToken generates base64-encoded UUID access and refresh tokens
func (a *AccessTokenGenDefault) GenerateAccessToken(data *AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	token := uuid.NewRandom()
	accesstoken = removePadding(base64.URLEncoding.EncodeToString([]byte(token)))

	if generaterefresh {
		rtoken := uuid.NewRandom()
		refreshtoken = removePadding(base64.URLEncoding.EncodeToString([]byte(rtoken)))
	}
	return
}
