package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/ory-am/hydra/client"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strconv"
)

type client struct {
	ep    string
	token *oauth2.Token
}

func New(endpoint string, token *oauth2.Token) Client {
	return &client{
		ep:    endpoint,
		token: token,
	}
}

var isAllowed struct {
	Allowed bool `json:"allowed"`
}

func (c *client) IsAllowed(ar *AuthorizeRequest) (bool, error) {
	request := gorequest.New()
	resp, body, errors := request.Post(c.ep+"/guard/allowed").Set("Authorization", c.token.Type()+" "+c.token.AccessToken).Send(ar).End()
	if len(errors) > 0 {
		return false, fmt.Errorf("Got errors: %v", errors)
	} else if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Got status code %s", resp.StatusCode)
	}

	if err := json.Unmarshal([]byte(body), &isAllowed); err != nil {
		return false, err
	}
	return isAllowed.Allowed, nil
}

func (c *client) IsAuthenticated(token string) (bool, error) {
	data := url.Values{}
	data.Set("token", token)

	client := &http.Client{}
	r, err := http.NewRequest("POST", c.ep+"/oauth2/introspect", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return false, err
	}
	r.Header.Add("Authorization", c.token.Type()+" "+c.token.AccessToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Got status code %s", resp.StatusCode)
	}

	var introspect struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&introspect); err != nil {
		return false, err
	}
	return introspect.Active, nil
}
