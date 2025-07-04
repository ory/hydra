// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type PageToken = interface {
	Parse(string) map[string]string
	Encode() string
}

var _ PageToken = new(StringPageToken)
var _ PageToken = new(MapPageToken)

type StringPageToken string

func (s StringPageToken) Parse(idField string) map[string]string {
	return map[string]string{idField: string(s)}
}

func (s StringPageToken) Encode() string {
	return string(s)
}

func NewStringPageToken(s string) (PageToken, error) {
	return StringPageToken(s), nil
}

type MapPageToken map[string]string

func (m MapPageToken) Parse(_ string) map[string]string {
	return map[string]string(m)
}

const pageTokenColumnDelim = "/"

func (m MapPageToken) Encode() string {
	elems := make([]string, 0, len(m))
	for k, v := range m {
		elems = append(elems, fmt.Sprintf("%s=%s", k, v))
	}

	// For now: use Base64 instead of URL escaping, as the Timestamp format we need to use can contain a `+` sign,
	// which represents a space in URLs, so it's not properly encoded by the Go library.
	return base64.RawStdEncoding.EncodeToString([]byte(strings.Join(elems, pageTokenColumnDelim)))
}

func NewMapPageToken(s string) (PageToken, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(string(b), pageTokenColumnDelim)

	r := map[string]string{}

	for _, p := range tokens {
		if columnName, value, found := strings.Cut(p, "="); found {
			r[columnName] = value
		}
	}

	return MapPageToken(r), nil
}

var _ PageTokenConstructor = NewMapPageToken
var _ PageTokenConstructor = NewStringPageToken

type PageTokenConstructor = func(string) (PageToken, error)
