// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"testing"
	"time"

	cristaljwt "github.com/cristalhq/jwt/v4"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
)

func ptr(d time.Duration) *time.Duration {
	return &d
}

var TestLifespans fosite.ClientLifespanConfig = fosite.ClientLifespanConfig{
	AuthorizationCodeGrantAccessTokenLifespan:  ptr(31 * time.Hour),
	AuthorizationCodeGrantIDTokenLifespan:      ptr(32 * time.Hour),
	AuthorizationCodeGrantRefreshTokenLifespan: ptr(33 * time.Hour),
	ClientCredentialsGrantAccessTokenLifespan:  ptr(34 * time.Hour),
	ImplicitGrantAccessTokenLifespan:           ptr(35 * time.Hour),
	ImplicitGrantIDTokenLifespan:               ptr(36 * time.Hour),
	JwtBearerGrantAccessTokenLifespan:          ptr(37 * time.Hour),
	PasswordGrantAccessTokenLifespan:           ptr(38 * time.Hour),
	PasswordGrantRefreshTokenLifespan:          ptr(39 * time.Hour),
	RefreshTokenGrantIDTokenLifespan:           ptr(40 * time.Hour),
	RefreshTokenGrantAccessTokenLifespan:       ptr(41 * time.Hour),
	RefreshTokenGrantRefreshTokenLifespan:      ptr(42 * time.Hour),
}

func RequireEqualDuration(t *testing.T, expected time.Duration, actual time.Duration, precision time.Duration) {
	delta := expected - actual
	if delta < 0 {
		delta = -delta
	}
	require.Less(t, delta, precision, fmt.Sprintf("expected %s; got %s", expected, actual))
}

func RequireEqualTime(t *testing.T, expected time.Time, actual time.Time, precision time.Duration) {
	delta := expected.Sub(actual)
	if delta < 0 {
		delta = -delta
	}
	require.Less(t, delta, precision, fmt.Sprintf(
		"expected %s; got %s",
		expected.Format(time.RFC3339Nano),
		actual.Format(time.RFC3339Nano),
	))
}

func ExtractJwtExpClaim(t *testing.T, token string) *time.Time {
	jwt, err := cristaljwt.ParseNoVerify([]byte(token))
	require.NoError(t, err)
	claims := &cristaljwt.RegisteredClaims{}
	require.NoError(t, json.Unmarshal(jwt.Claims(), claims))
	if claims.ExpiresAt == nil {
		return nil
	}
	return &claims.ExpiresAt.Time
}

func ParseFormPostResponse(t *testing.T, redirectURL string, resp io.Reader) (authorizationCode, stateFromServer, iDToken string, token goauth.Token, customParameters url.Values, rFC6749Error map[string]string) {
	token = goauth.Token{}
	rFC6749Error = map[string]string{}
	customParameters = url.Values{}

	doc, err := html.Parse(resp)
	require.NoError(t, err)

	//doc>html>body
	body := findBody(doc.FirstChild.FirstChild)
	require.Equal(t, "body", body.Data)

	htmlEvent := body.Attr[0].Key
	require.Equal(t, "onload", htmlEvent)

	onLoadFunc := body.Attr[0].Val
	require.Equal(t, "javascript:document.forms[0].submit()", onLoadFunc)

	form := getNextNoneTextNode(body.FirstChild)
	require.NotNil(t, form)
	require.Equal(t, "form", form.Data)

	for _, attr := range form.Attr {
		if attr.Key == "method" {
			require.Equal(t, "post", attr.Val)
		} else {
			require.Equal(t, redirectURL, attr.Val)
		}
	}

	for node := getNextNoneTextNode(form.FirstChild); node != nil; node = getNextNoneTextNode(node.NextSibling) {
		var k, v string
		for _, attr := range node.Attr {
			if attr.Key == "name" {
				k = attr.Val
			} else if attr.Key == "value" {
				v = attr.Val
			}
		}

		switch k {
		case "state":
			stateFromServer = v
		case "code":
			authorizationCode = v
		case "expires_in":
			expires, err := strconv.Atoi(v)
			require.NoError(t, err)
			token.Expiry = time.Now().UTC().Add(time.Duration(expires) * time.Second)
		case "access_token":
			token.AccessToken = v
		case "token_type":
			token.TokenType = v
		case "refresh_token":
			token.RefreshToken = v
		case "error":
			rFC6749Error["ErrorField"] = v
		case "error_hint":
			rFC6749Error["HintField"] = v
		case "error_description":
			rFC6749Error["DescriptionField"] = v
		case "id_token":
			iDToken = v
		default:
			customParameters.Add(k, v)
		}
	}

	return
}

func getNextNoneTextNode(node *html.Node) *html.Node {
	nextNode := node.NextSibling
	if nextNode != nil && nextNode.Type == html.TextNode {
		nextNode = getNextNoneTextNode(node.NextSibling)
	}

	return nextNode
}

func findBody(node *html.Node) *html.Node {
	if node != nil {
		if node.Data == "body" {
			return node
		}
		return findBody(node.NextSibling)
	}

	return nil
}
