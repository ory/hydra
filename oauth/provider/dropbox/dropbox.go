package dropbox

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/oauth/provider"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io/ioutil"
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
		Given       string `json:"given_name,omitempty"`
		Surname     string `json:"surname,omitempty"`
		FamilyName  string `json:"familiar_name,omitempty"`
		DisplayName string `json:"display_name,omitempty"`
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
				TokenURL: "https://www.dropbox.com/1/oauth2/token",
			},
		},
	}
}

func (d *dropbox) GetAuthCodeURL(state string) string {
	return d.conf.AuthCodeURL(state)
}

func (d *dropbox) Exchange(code string) (Session, error) {
	conf := *d.conf
	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.Errorf("Token is not valid: %v", token)
	}

	c := conf.Client(ctx, token)
	rawurl := fmt.Sprintf("%s/%s?%s", d.api, "users/get_current_account", nil)
	response, err := c.Get(rawurl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var acc Account
	if err = json.Unmarshal(body, &acc); err != nil {
		return nil, err
	}

	return &DefaultSession{
		RemoteSubject: acc.ID,
		Extra:         acc,
		Token:         token,
	}, nil
}

func (d *dropbox) GetID() string {
	return d.id
}
