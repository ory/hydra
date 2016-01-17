package dropbox

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/oauth/provider"
	"golang.org/x/oauth2"
	"net/http"
)

type dropbox struct {
	id    string
	conf  *oauth2.Config
	token *oauth2.Token
	api   string
}

type Account struct {
	ID          string                 `json:"account_id"`
	Email       string                 `json:"email"`
	Locale      string                 `json:"locale"`
	ReferralURL string                 `json:"referral_link"`
	IsPaired    bool                   `json:"is_paired"`
	Type        map[string]interface{} `json:"account_type"`
	Name        struct {
		Given        string `json:"given_name,omitempty"`
		Surname      string `json:"surname,omitempty"`
		FamiliarName string `json:"familiar_name,omitempty"`
		DisplayName  string `json:"display_name,omitempty"`
	} `json:"name"`
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

	var acc Account
	if err = json.NewDecoder(response.Body).Decode(&acc); err != nil {
		return nil, err
	}

	return &DefaultSession{
		RemoteSubject: acc.ID,
		Extra: map[string]interface{}{
			"account_id":    acc.ID,
			"email":         acc.Email,
			"locale":        acc.Locale,
			"referral_link": acc.ReferralURL,
			"is_paired":     acc.IsPaired,
			"account_type":  acc.Type,
			"name": map[string]interface{}{
				"given_name":    acc.Name.Given,
				"surname":       acc.Name.Surname,
				"familiar_name": acc.Name.FamiliarName,
				"display_name":  acc.Name.DisplayName,
			},
		},
	}, nil
}

func (d *dropbox) GetID() string {
	return d.id
}
