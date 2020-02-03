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

	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/urlx"

	hydra "github.com/ory/hydra/internal/httpclient/client"
)

var hydraURL = urlx.ParseOrPanic(os.Getenv("HYDRA_ADMIN_URL"))
var client = hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{hydraURL.Scheme}, Host: hydraURL.Host, BasePath: hydraURL.Path})

func login(rw http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("login_challenge")
	res, err := client.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(challenge))
	if err != nil {
		log.Fatalf("Unable to fetch clogin request: %s", err)
	}
	lr := res.Payload

	var v *models.CompletedRequest
	if strings.Contains(lr.RequestURL, "mockLogin=accept") {
		remember := false
		if strings.Contains(lr.RequestURL, "rememberLogin=yes") {
			remember = true
		}

		var vr *admin.AcceptLoginRequestOK
		vr, err = client.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(challenge).
			WithBody(&models.AcceptLoginRequest{
				Subject:  pointerx.String("the-subject"),
				Remember: remember,
			}))
		if vr != nil {
			v = vr.Payload
		}
	} else {
		var vr *admin.RejectLoginRequestOK
		vr, err = client.Admin.RejectLoginRequest(admin.NewRejectLoginRequestParams().
			WithLoginChallenge(challenge).
			WithBody(&models.RejectRequest{
				Error: "invalid_request",
			}))
		if vr != nil {
			v = vr.Payload
		}
	}
	if err != nil {
		log.Fatalf("Unable to accept/reject login request: %s", err)
	}
	http.Redirect(rw, r, v.RedirectTo, http.StatusFound)
}

func consent(rw http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("consent_challenge")

	rrr, err := client.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(challenge))
	if err != nil {
		log.Fatalf("Unable to fetch consent request: %s", err)
	}
	o := rrr.Payload

	var v *models.CompletedRequest
	if strings.Contains(o.RequestURL, "mockConsent=accept") {
		remember := false
		if strings.Contains(o.RequestURL, "rememberConsent=yes") {
			remember = true
		}
		value := "bar"
		if o.Skip {
			value = "rab"
		}

		var vr *admin.AcceptConsentRequestOK
		vr, err = client.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
			WithConsentChallenge(challenge).
			WithBody(&models.AcceptConsentRequest{
				GrantScope: o.RequestedScope,
				Remember:   remember,
				Session: &models.ConsentRequestSession{
					AccessToken: map[string]interface{}{"foo": value},
					IDToken:     map[string]interface{}{"baz": value},
				},
			}))
		if vr != nil {
			v = vr.Payload
		}
	} else {
		var vr *admin.RejectConsentRequestOK
		vr, err = client.Admin.RejectConsentRequest(
			admin.NewRejectConsentRequestParams().WithConsentChallenge(challenge).
				WithBody(
					&models.RejectRequest{
						Error: "invalid_request",
					}),
		)
		if vr != nil {
			v = vr.Payload
		}
	}
	if err != nil {
		log.Fatalf("Unable to accept/reject consent request: %s", err)
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
