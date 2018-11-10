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
	"strings"

	"github.com/ory/hydra/sdk/go/hydra/swagger"
)

var client = swagger.NewAdminApiWithBasePath(os.Getenv("HYDRA_ADMIN_URL"))

func login(rw http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("login_challenge")
	lr, resp, err := client.GetLoginRequest(challenge)
	if err != nil {
		log.Fatalf("Unable to fetch clogin request: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch login request, got status code %d", resp.StatusCode)
	}

	var v *swagger.CompletedRequest
	if strings.Contains(lr.RequestUrl, "mockLogin=accept") {
		remember := false
		if strings.Contains(lr.RequestUrl, "rememberLogin=yes") {
			remember = true
		}
		v, resp, err = client.AcceptLoginRequest(challenge, swagger.AcceptLoginRequest{
			Subject:  "the-subject",
			Remember: remember,
		})
	} else {
		v, resp, err = client.RejectLoginRequest(challenge, swagger.RejectRequest{
			Error_: "invalid_request",
		})
	}
	if err != nil {
		log.Fatalf("Unable to accept/reject login request: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to accept/reject login request, got status code %d", resp.StatusCode)
	}
	http.Redirect(rw, r, v.RedirectTo, http.StatusFound)
}

func consent(rw http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("consent_challenge")
	o, resp, err := client.GetConsentRequest(challenge)
	if err != nil {
		log.Fatalf("Unable to fetch consent request: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch consent request, got status code %d", resp.StatusCode)
	}

	var v *swagger.CompletedRequest
	if strings.Contains(o.RequestUrl, "mockConsent=accept") {
		remember := false
		if strings.Contains(o.RequestUrl, "rememberConsent=yes") {
			remember = true
		}
		value := "bar"
		if o.Skip == true {
			value = "rab"
		}

		v, resp, err = client.AcceptConsentRequest(challenge, swagger.AcceptConsentRequest{
			GrantScope: o.RequestedScope,
			Remember:   remember,
			Session: swagger.ConsentRequestSession{
				AccessToken: map[string]interface{}{"foo": value},
				IdToken:     map[string]interface{}{"baz": value},
			},
		})
	} else {
		v, resp, err = client.RejectConsentRequest(challenge, swagger.RejectRequest{
			Error_: "invalid_request",
		})
	}
	if err != nil {
		log.Fatalf("Unable to accept/reject consent request: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to accept/reject consent request, got status code %d", resp.StatusCode)
	}
	http.Redirect(rw, r, v.RedirectTo, http.StatusFound)
}

func errh(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, r.URL.Query().Get("error")+" "+r.URL.Query().Get("error_debug"), http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/consent", consent)
	http.HandleFunc("/error", errh)
	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
