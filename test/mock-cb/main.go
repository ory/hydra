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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func callback(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("error") != "" {
		http.Error(rw, "error happened in callback: "+r.URL.Query().Get("error")+" "+r.URL.Query().Get("error_description")+" "+r.URL.Query().Get("error_debug"), http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	conf := oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  strings.TrimRight(os.Getenv("HYDRA_URL"), "/") + "/oauth2/auth",
			TokenURL: strings.TrimRight(os.Getenv("HYDRA_URL"), "/") + "/oauth2/token",
		},
		RedirectURL: os.Getenv("REDIRECT_URL"),
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(&struct {
		IDToken      string    `json:"id_token"`
		AccessToken  string    `json:"access_token"`
		TokenType    string    `json:"token_type,omitempty"`
		RefreshToken string    `json:"refresh_token,omitempty"`
		Expiry       time.Time `json:"expiry,omitempty"`
	}{
		IDToken:      fmt.Sprintf("%s", token.Extra("id_token")),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/callback", callback)
	port := "4445"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Fatal(http.ListenAndServe(":"+port, nil)) // #nosec G114
}
