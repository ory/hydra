// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	hydra "github.com/ory/hydra-client-go/v2"

	"golang.org/x/oauth2"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/urlx"
)

var hydraURL = urlx.ParseOrPanic(os.Getenv("HYDRA_ADMIN_URL"))
var sdk = hydra.NewAPIClient(hydra.NewConfiguration())

func init() {
	sdk.GetConfig().Servers = hydra.ServerConfigurations{{URL: hydraURL.String()}}
}

type oauth2token struct {
	IDToken      string    `json:"id_token"`
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

var printToken, printCookie bool

func init() {
	flag.BoolVar(&printToken, "print-token", false, "")
	flag.BoolVar(&printCookie, "print-cookie", false, "")
}

func main() {
	flag.Parse()
	conf := oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  strings.TrimRight(os.Getenv("HYDRA_URL"), "/") + "/oauth2/auth",
			TokenURL: strings.TrimRight(os.Getenv("HYDRA_URL"), "/") + "/oauth2/token",
		},
		Scopes:      strings.Split(os.Getenv("OAUTH2_SCOPE"), ","),
		RedirectURL: os.Getenv("REDIRECT_URL"),
	}
	au := conf.AuthCodeURL("some-stupid-state-foo") + os.Getenv("OAUTH2_EXTRA")
	c, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatalf("Unable to create cookie jar: %s", err)
	}

	u, _ := url.Parse("http://127.0.0.1")
	if os.Getenv("AUTH_COOKIE") != "" {
		c.SetCookies(u, []*http.Cookie{{Name: "oauth2_authentication_session", Value: os.Getenv("AUTH_COOKIE")}})
	}

	resp, err := (&http.Client{
		Jar: c,
		// Hack to fix cookie across domains
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 && req.Header.Get("cookie") == "" {
				req.Header.Set("Cookie", via[len(via)-1].Header.Get("Cookie"))
			}

			return nil
		},
	}).Get(au)
	cmdx.CheckResponse(err, http.StatusOK, resp)
	defer resp.Body.Close() //nolint:errcheck

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read body: %s", err)
	}

	for _, c := range c.Cookies(u) {
		if c.Name == "oauth2_authentication_session" {
			if printCookie {
				fmt.Print(c.Value)
			}
		}
	}

	var token oauth2token
	if err := json.Unmarshal(out, &token); err != nil {
		log.Fatalf("Unable transform to token: %s", err)
	}

	checkTokenResponse(token)
	for i := 0; i <= 5; i++ {
		token = refreshToken(token)
		checkTokenResponse(token)
	}

	newToken := refreshToken(token)
	if printToken {
		fmt.Printf("%s", newToken.AccessToken)
	}

	// refreshing the same token twice does not work
	resp, err = refreshTokenRequest(token)
	cmdx.CheckResponse(err, http.StatusBadRequest, resp)
	defer resp.Body.Close() //nolint:errcheck
}

func refreshTokenRequest(token oauth2token) (*http.Response, error) {
	req, err := http.NewRequest("POST", strings.TrimRight(os.Getenv("HYDRA_URL"), "/")+"/oauth2/token", bytes.NewBufferString(url.Values{
		"refresh_token": {token.RefreshToken},
		"grant_type":    {"refresh_token"},
	}.Encode()))
	cmdx.Must(err, "%s", err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv("OAUTH2_CLIENT_ID"), os.Getenv("OAUTH2_CLIENT_SECRET"))
	return http.DefaultClient.Do(req)
}

func refreshToken(token oauth2token) (result oauth2token) {
	resp, err := refreshTokenRequest(token)
	cmdx.CheckResponse(err, http.StatusOK, resp)
	defer resp.Body.Close() //nolint:errcheck

	err = json.NewDecoder(resp.Body).Decode(&result)
	cmdx.Must(err, "Unable to decode refresh token: %s", err)
	return result
}

func checkTokenResponse(token oauth2token) {
	if token.RefreshToken == "" {
		log.Fatalf("Expected a refresh token but none received: %+v", token)
	}

	// This value oscillates between bar and rab, depending on whether authorization was remembered or not. Check
	// mock-lcp which sets the value
	expectedValue := "bar"
	if strings.Contains(os.Getenv("OAUTH2_EXTRA"), "prompt=none") {
		expectedValue = "rab"
	}

	if os.Getenv("OAUTH2_ACCESS_TOKEN_STRATEGY") == "jwt" {
		parts := strings.Split(token.AccessToken, ".")

		if len(parts) != 3 {
			log.Fatalf("JWT Access Token does not seem to have three parts: %d - %+v - %v", len(parts), token, parts)
		}

		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			log.Fatalf("Unable to decode id token segment: %s", err)
		}

		var claims map[string]interface{}
		if err := json.Unmarshal(payload, &claims); err != nil {
			log.Fatalf("Unable to unmarshal id token body: %s", err)
		}

		if fmt.Sprintf("%s", claims["sub"]) != "the-subject" {
			log.Fatalf("Expected subject from access token to be %s but got %s", "the-subject", claims["sub"])
		}

		ext := claims["ext"].(map[string]interface{})
		if ext["foo"] != expectedValue {
			log.Fatalf("Expected extra field \"foo\" from access token to be \"%s\" but got %s", expectedValue, ext["foo"])
		}
	}

	intro, resp, err := sdk.OAuth2API.IntrospectOAuth2Token(context.Background()).Token(token.AccessToken).Execute()
	defer resp.Body.Close() //nolint:errcheck
	if err != nil {
		log.Fatalf("Unable to introspect OAuth2 token: %s", err)
	}

	if !intro.Active {
		log.Fatalf("Expected token to be active: %s", token.AccessToken)
	}

	if *intro.Sub != "the-subject" {
		log.Fatalf("Expected subject from access token to be %s but got %s", "the-subject", *intro.Sub)
	}

	if intro.Ext["foo"] != expectedValue {
		log.Fatalf("Expected extra field \"foo\" from access token to be \"%s\" but got %s", expectedValue, intro.Ext["foo"])
	}

	idt := token.IDToken
	if len(idt) == 0 {
		log.Fatalf("ID Token does not seem to be set: %+v", token)
	}

	parts := strings.Split(idt, ".")
	if len(parts) != 3 {
		log.Fatalf("ID Token does not seem to have three parts: %d - %+v - %v", len(parts), token, parts)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatalf("Unable to decode id token segment: %s", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		log.Fatalf("Unable to unmarshal id token body: %s", err)
	}

	if fmt.Sprintf("%s", claims["sub"]) != "the-subject" {
		log.Fatalf("Expected subject from id token to be %s but got %s", "the-subject", claims["sub"])
	}

	if fmt.Sprintf("%s", claims["baz"]) != expectedValue {
		log.Fatalf("Expected extra field \"baz\" from access token to be \"%s\" but got \"%s\"", expectedValue, claims["baz"])
	}
}
