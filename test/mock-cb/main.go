// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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
