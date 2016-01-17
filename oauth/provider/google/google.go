package google

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/oauth/provider"
	"golang.org/x/oauth2"
	gauth "golang.org/x/oauth2/google"
	"net/http"
)

type google struct {
	id   string
	api  string
	conf *oauth2.Config
}

type claims struct {
	Issuer     string `json:"iss"`
	Subject    string `json:"sub"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Locale     string `json:"locale"`
}

func New(id, client, secret, redirectURL string) *google {
	return &google{
		id:  id,
		api: "https://www.googleapis.com",
		conf: &oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			Scopes:       []string{"openid", "email", "profile"},
			RedirectURL:  redirectURL,
			Endpoint:     gauth.Endpoint,
		},
	}
}

func (d *google) GetAuthenticationURL(state string) string {
	return d.conf.AuthCodeURL(state)
}

func (d *google) FetchSession(code string) (Session, error) {
	conf := *d.conf
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.Errorf("Token is not valid: %v", token)
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.Errorf("Token is not valid: %v", idToken)
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", d.api, "oauth2/v3/tokeninfo?id_token="+idToken))
	if err != nil {
		return nil, errors.Errorf("Could not validate id token because %s", err)
	}
	defer resp.Body.Close()

	var profile claims
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, errors.Errorf("Could not validate id token because %s", err)
	}

	return &DefaultSession{
		RemoteSubject: profile.Subject,
		Extra: map[string]interface{}{
			"iss":         profile.Issuer,
			"sub":         profile.Subject,
			"email":       profile.Email,
			"picture":     profile.Picture,
			"locale":      profile.Locale,
			"given_name":  profile.GivenName,
			"name":        profile.Name,
			"family_name": profile.FamilyName,
		},
	}, nil
}

func (d *google) GetID() string {
	return d.id
}
