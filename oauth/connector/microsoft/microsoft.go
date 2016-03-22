package microsoft

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/oauth/connector"
	"golang.org/x/oauth2"
	"net/http"
)

// Read up on: https://dev.onedrive.com/auth/msa_oauth.htm

type microsoft struct {
	id    string
	conf  *oauth2.Config
	token *oauth2.Token
	api   string
}

func New(id, client, secret, redirectURL string) *microsoft {
	return &microsoft{
		id:  id,
		api: "https://apis.live.net",
		conf: &oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"wl.signin", "wl.emails"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://login.live.com/oauth20_authorize.srf",
				TokenURL: "https://login.live.com/oauth20_token.srf",
			},
		},
	}
}

func (d *microsoft) GetAuthenticationURL(state string) string {
	return d.conf.AuthCodeURL(state)
}

func (d *microsoft) FetchSession(code string) (Session, error) {
	conf := *d.conf
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.Errorf("Token is not valid: %v", token)
	}

	c := conf.Client(oauth2.NoContext, token)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", d.api, "v5.0/me"), nil)
	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Could not fetch account data because %s", err)
	}

	var acc map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&acc); err != nil {
		return nil, err
	}

	return &DefaultSession{
		RemoteSubject: fmt.Sprintf("%s", acc["id"]),
		Extra:         acc,
	}, nil
}

func (d *microsoft) GetID() string {
	return d.id
}
