package slack

import (
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/oauth/provider"
	"golang.org/x/oauth2"
	"net/http"
	"encoding/json"
)

type facebook struct {
	id    string
	conf  *oauth2.Config
	token *oauth2.Token
	api   string
}

func New(id, client, secret, redirectURL string) *facebook {
	return &facebook{
		id:  id,
		api: "https://slack.com/api",
		conf: &oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			RedirectURL:  redirectURL,
			Scopes: []string{"identity.basic", "identity.email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://slack.com/oauth/authorize",
				TokenURL: "https://slack.com/api/oauth.access",
			},
		},
	}
}

func (d *facebook) GetAuthenticationURL(state string) string {
	return d.conf.AuthCodeURL(state)
}

func (d *facebook) FetchSession(code string) (Session, error) {
	conf := *d.conf
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.Errorf("Token is not valid: %v", token)
	}


	c := conf.Client(oauth2.NoContext, token)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?token=%s", d.api, "users.identity", token.AccessToken), nil)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Could not fetch account data because got status code %d", resp.StatusCode)
	}

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, errors.Errorf("Could not validate id token because %s", err)
	}

	return &DefaultSession{
		RemoteSubject: fmt.Sprintf("%s", profile["id"]),
		Extra:         profile,
	}, nil
}

func (d *facebook) GetID() string {
	return d.id
}
