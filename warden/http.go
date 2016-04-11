package warden

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	"github.com/ory-am/common/pkg"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"github.com/ory-am/hydra/handler"
)

var isAllowed struct {
	Allowed bool `json:"allowed"`
}

type HTTPWarden struct {
	Client   *http.Client

	Endpoint url.URL
}

func (w *HTTPWarden) SetClient(c *clientcredentials.Config) {
	w.Client = c.Client(oauth2.NoContext)
}

func (w *HTTPWarden) ActionAllowed(token string, action *Action) (*Context, error) {
	json.New
	req, err := w.Client.Post(w.Endpoint.String(), "application/json", action)
	if err != nil {
		return nil, errors.New(err)
	}

	var decision handler.DecisionResponse
	if err := json.NewDecoder(req).Decode(&decision); err != nil {
		return nil, errors.New(err)
	}

	ctx := &Context{

	}

	return nil, nil
}

func (c *HTTPWarden) Authorized(token string, scopes ...string) (*Context, error) {

}

func (c *HTTPWarden) HTTPAuthorized(r *http.Request, scopes ...string) (*Context, error) {

}

func (c *HTTPWarden) HTTPActionAllowed(r *http.Request, scopes ...string) (*Context, error) {

}



func (c *HTTPWarden) SetClientToken(token *oauth2.Token) {
	c.clientToken = token
}

func (c *HTTPWarden) IsRequestAllowed(req *http.Request, resource, permission, owner string) (bool, error) {
	var token *osin.BearerAuth
	if token = osin.CheckBearerAuth(req); token == nil {
		return false, errors.New("No token given.")
	} else if token.Code == "" {
		return false, errors.New("No token given.")
	}
	env := middleware.NewEnv(req)
	env.Owner(owner)
	return c.IsAllowed(&Action{Token: token.Code, Resource: resource, Permission: permission, Context: env.Ctx()})
}

func (c *HTTPWarden) IsAllowed(ar *Action) (bool, error) {
	return isValidAuthorizeRequest(c, ar, true)
}

func (c *HTTPWarden) IsAuthenticated(token string) (bool, error) {
	return isValidAuthenticationRequest(c, token, true)
}

func isValidAuthenticationRequest(c *HTTPWarden, token string, retry bool) (bool, error) {
	data := url.Values{}
	data.Set("token", token)
	request := gorequest.New()
	resp, body, errs := request.Post(pkg.JoinURL(c.ep, "/oauth2/introspect")).Type("form").SetBasicAuth(c.ClientConfig.ClientID, c.ClientConfig.ClientSecret).Set("Connection", "close").SendString(data.Encode()).End()
	if len(errs) > 0 {
		return false, errors.Errorf("Got errors: %v", errs)
	} else if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("Status code %d is not 200: %s", resp.StatusCode, body)
	}

	if retry && resp.StatusCode == http.StatusUnauthorized {
		var err error
		if c.clientToken, err = c.ClientConfig.Token(oauth2.NoContext); err != nil {
			return false, errors.New(err)
		} else if c.clientToken == nil {
			return false, errors.New("Access token could not be retrieved")
		}
		return isValidAuthenticationRequest(c, token, false)
	} else if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Status code %d is not 200", resp.StatusCode)
	}

	var introspect struct {
		Active bool `json:"active"`
	}

	if err := json.Unmarshal([]byte(body), &introspect); err != nil {
		return false, err
	} else if !introspect.Active {
		return false, errors.New("Authentication denied")
	}
	return introspect.Active, nil
}

func isValidAuthorizeRequest(c *HTTPWarden, ar *Action, retry bool) (bool, error) {
	request := gorequest.New()
	resp, body, errs := request.Post(pkg.JoinURL(c.ep, "/guard/allowed")).SetBasicAuth(c.ClientConfig.ClientID, c.ClientConfig.ClientSecret).Set("Content-Type", "application/json").Set("Connection", "close").Send(*ar).End()
	if len(errs) > 0 {
		return false, errors.Errorf("Got errors: %v", errs)
	} else if retry && resp.StatusCode == http.StatusUnauthorized {
		var err error
		if c.clientToken, err = c.ClientConfig.Token(oauth2.NoContext); err != nil {
			return false, errors.New(err)
		} else if c.clientToken == nil {
			return false, errors.New("Access token could not be retrieved")
		}
		return isValidAuthorizeRequest(c, ar, false)
	} else if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("Status code %d is not 200: %s", resp.StatusCode, body)
	}

	if err := json.Unmarshal([]byte(body), &isAllowed); err != nil {
		return false, errors.Errorf("Could not unmarshall body because %s", err.Error())
	}

	if !isAllowed.Allowed {
		return false, errors.New("Authroization denied")
	}
	return isAllowed.Allowed, nil
}
