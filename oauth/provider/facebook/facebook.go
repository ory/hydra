package facebook

import (
	"fmt"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/oauth/provider"
	"golang.org/x/oauth2"
	fb "github.com/huandu/facebook"
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
		api: "https://graph.facebook.com",
		conf: &oauth2.Config{
			ClientID:     client,
			ClientSecret: secret,
			RedirectURL:  redirectURL,
			Scopes: []string{"email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/dialog/oauth",
				TokenURL: "https://graph.facebook.com/v2.3/oauth/access_token",
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

	session := &fb.Session{HttpClient: conf.Client(oauth2.NoContext, token)}
	res, err := session.Get("/me", fb.Params{
		"fields": []string{"email", "id", "first_name", "last_name"},
	})
	if err != nil {
		return nil, err
	}

	var acc map[string]interface{}
	if err = res.Decode(&acc); err != nil {
		return nil, err
	}

	return &DefaultSession{
		RemoteSubject: fmt.Sprintf("%s", res.Get("id")),
		Extra:         acc,
	}, nil
}

func (d *facebook) GetID() string {
	return d.id
}
