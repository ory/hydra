package dropbox

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/endpoint/connector"
	"golang.org/x/oauth2"
	"net/http"
)

type dropbox struct {
	id    string
	conf  *oauth2.Config
	token *oauth2.Token
	api   string
}

func New(id, client, secret, redirectURL string) *dropbox {
	return &dropbox{
		id:  id,
		api: "https://api.dropbox.com/2",
		conf: &oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			RedirectURL:  redirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.dropbox.com/1/oauth2/authorize",
				TokenURL: "https://api.dropbox.com/1/oauth2/token",
			},
		},
	}
}

func (d *dropbox) GetAuthenticationURL(state string) string {
	return d.conf.AuthCodeURL(state)
}

func (d *dropbox) FetchSession(code string) (Session, error) {
	conf := *d.conf
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.Errorf("Token is not valid: %v", token)
	}

	c := conf.Client(oauth2.NoContext, token)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", d.api, "users/get_current_account"), nil)
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
		RemoteSubject: fmt.Sprintf("%s", acc["account_id"]),
		Extra:         acc,
	}, nil
}

func (d *dropbox) GetID() string {
	return d.id
}
