// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type IntrospectForm struct {
	Token  string
	Scopes []string
}

type IntrospectResponse struct {
	Active    bool     `json:"active"`
	ClientID  string   `json:"client_id,omitempty"`
	Scope     string   `json:"scope,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Username  string   `json:"username,omitempty"`
}

type Introspect struct {
	endpointURL string
	client      *http.Client
}

func (c *Introspect) IntrospectToken(
	ctx context.Context,
	form IntrospectForm,
	header map[string]string,
) (*IntrospectResponse, error) {
	data := url.Values{}
	data.Set("token", form.Token)
	data.Set("scope", strings.Join(form.Scopes, " "))

	request, err := c.getRequest(ctx, data, header)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if c := response.StatusCode; c < 200 || c > 299 {
		return nil, &RequestError{
			Response: response,
			Body:     body,
		}
	}

	result := &IntrospectResponse{}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Introspect) getRequest(
	ctx context.Context,
	data url.Values,
	header map[string]string,
) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", c.endpointURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for header, value := range header {
		request.Header.Set(header, value)
	}

	return request, nil
}

func NewIntrospectClient(endpointURL string) *Introspect {
	return &Introspect{
		endpointURL: endpointURL,
		client:      &http.Client{},
	}
}
