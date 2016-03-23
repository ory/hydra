package google

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/endpoint/connector"
	"golang.org/x/oauth2"
	gauth "golang.org/x/oauth2/google"
	"net/http"
)

type google struct {
	id   string
	api  string
	conf *oauth2.Config
}

func New(id, client, secret, redirectURL string) *google {
	return &google{
		id:  id,
		api: "https://www.googleapis.com",
		conf: &oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			Scopes: []string{
				"email",
				"profile",
				"https://www.googleapis.com/auth/plus.login",
				"https://www.googleapis.com/auth/plus.me",
			},
			RedirectURL: redirectURL,
			Endpoint:    gauth.Endpoint,
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

	c := conf.Client(oauth2.NoContext, token)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", d.api, "plus/v1/people/me"), nil)
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

func (d *google) GetID() string {
	return d.id
}
