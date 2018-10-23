/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"

	"github.com/ory/hydra/sdk/go/hydra/swagger"
)

func main() {
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
	if err != nil {
		log.Fatalf("Unable to make request: %s", err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read body: %s", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Got status code %d and body %s", resp.StatusCode, out)
	}

	var token oauth2.Token
	if err := json.Unmarshal(out, &token); err != nil {
		log.Fatalf("Unable transform to token: %s", err)
	}

	for _, c := range c.Cookies(u) {
		if c.Name == "oauth2_authentication_session" {
			fmt.Print(c.Value)
		}
	}

	resp, err = http.PostForm(strings.TrimRight(os.Getenv("HYDRA_ADMIN_URL"), "/")+"/oauth2/introspect", url.Values{"token": {token.AccessToken}})
	if err != nil {
		log.Fatalf("Unable to make introspection request: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to make introspection request: got status code %d", resp.StatusCode)
	}

	var intro swagger.OAuth2TokenIntrospection
	if err := json.NewDecoder(resp.Body).Decode(&intro); err != nil {
		log.Fatalf("Unable to decode introspection response: %s", err)
	}
	resp.Body.Close()

	if intro.Sub != "the-subject" {
		log.Fatalf("Expected subject from access token to be %s but got %s", "the-subject", intro.Sub)
	}

	if intro.Ext["foo"] != "bar" {
		log.Fatalf("Expected extra field \"foo\" from access token to be \"bar\" but got %s", intro.Ext["foo"])
	}

	payload, err := jwt.DecodeSegment(strings.Split(fmt.Sprintf("%s", token.Extra("id_token")), ".")[1])
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

	if fmt.Sprintf("%s", claims["foo"]) != "bar" {
		log.Fatalf("Expected extra field \"foo\" from access token to be \"bar\" but got %s", intro.Ext["foo"])
	}
}
