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
	"log"
	"net/http"
	"os"

	"io/ioutil"
	"net/http/cookiejar"
	"strings"

	"fmt"
	"net/url"

	"golang.org/x/oauth2"
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
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Got status code %d and body %s", resp.StatusCode, out)
	}

	for _, c := range c.Cookies(u) {
		if c.Name == "oauth2_authentication_session" {
			fmt.Print(c.Value)
		}
	}
}
