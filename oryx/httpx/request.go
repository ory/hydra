// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// NewRequestJSON returns a new JSON *http.Request.
func NewRequestJSON(method, url string, data interface{}) (*http.Request, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(data); err != nil {
		return nil, errors.WithStack(err)
	}
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// NewRequestForm returns a new POST Form *http.Request.
func NewRequestForm(method, url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// MustNewRequest returns a new *http.Request or fatals.
func MustNewRequest(method, url string, body io.Reader, contentType string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return req
}
