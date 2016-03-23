package signin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/endpoint/connector"
	"github.com/parnurzeal/gorequest"
)

type signin struct {
	id         string
	login      string
	redirectTo string
}

type payload struct {
	Subject string `json:"subject"`
}

func New(id, loginURL, redirectToURL string) *signin {
	return &signin{
		id:         id,
		login:      loginURL,
		redirectTo: redirectToURL,
	}
}

func (d *signin) GetAuthenticationURL(state string) string {
	// FIXME does not work if redirect contains query params
	return fmt.Sprintf("%s?state=%s&redirect_uri=%s", d.login, state, d.redirectTo)
}

func (d *signin) FetchSession(code string) (Session, error) {
	request := gorequest.New()

	// FIXME does not work if redirect contains query params
	resp, body, errs := request.Get(fmt.Sprintf("%s?verify=%s", d.login, code)).End()
	if len(errs) > 0 {
		return nil, errors.Errorf("Could not exchange code: %s", errs)
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Could not exchange code, received status: %d", resp.StatusCode)
	}

	var p payload
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		return nil, errors.Errorf("Could not parse answer: %v", err)
	}

	if p.Subject == "" {
		return nil, errors.Errorf("Field subject is empty, got %s", body)
	}

	return &DefaultSession{
		ForceLocalSubject: p.Subject,
		Extra:             map[string]interface{}{},
	}, nil
}

func (d *signin) GetID() string {
	return d.id
}
