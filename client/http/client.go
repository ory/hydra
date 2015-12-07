package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/middleware"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strconv"
)

var isAllowed struct {
	Allowed bool `json:"allowed"`
}

type HTTPClient struct {
	ep          string
	clientToken *oauth2.Token
}

func New(endpoint string, token *oauth2.Token) *HTTPClient {
	return &HTTPClient{
		ep:          endpoint,
		clientToken: token,
	}
}

func (c *HTTPClient) IsRequestAllowed(req *http.Request, resource, permission, owner string) (bool, error) {
	var token *osin.BearerAuth
	if token = osin.CheckBearerAuth(req); token == nil {
		return false, errors.New("No token given.")
	} else if token.Code == "" {
		return false, errors.New("No token given.")
	}
	env := middleware.NewEnv(req)
	env.Owner(owner)
	return c.IsAllowed(&AuthorizeRequest{Token: token.Code, Resource: resource, Permission: permission, Context: env.Ctx()})
}

func (c *HTTPClient) IsAllowed(ar *AuthorizeRequest) (bool, error) {
	request := gorequest.New()
	resp, body, errs := request.Post(c.ep+"/guard/allowed").Set("Authorization", c.clientToken.Type()+" "+c.clientToken.AccessToken).Set("Content-Type", "application/json").Send(*ar).End()
	if len(errs) > 0 {
		return false, errors.Errorf("Got errors: %v", errs)
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

func (c *HTTPClient) IsAuthenticated(token string) (bool, error) {
	data := url.Values{}
	client := &http.Client{}
	data.Set("token", token)
	r, err := http.NewRequest("POST", c.ep+"/oauth2/introspect", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return false, err
	}

	r.Header.Add("Authorization", c.clientToken.Type()+" "+c.clientToken.AccessToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Status code %d is not 200", resp.StatusCode)
	}

	var introspect struct {
		Active bool `json:"active"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&introspect); err != nil {
		return false, err
	}

	if !introspect.Active {
		return false, errors.New("Authentication denied")
	}
	return introspect.Active, nil
}
