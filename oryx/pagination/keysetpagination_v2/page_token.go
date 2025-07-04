// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/ssoready/hyrumtoken"

	"github.com/ory/herodot"
)

type (
	PageToken struct {
		testNow func() time.Time
		cols    []Column
	}
	jsonPageToken = struct {
		ExpiresAt time.Time `json:"e"`
		Cols      []Column  `json:"c"`
	}
	Column struct {
		Name  string `json:"n"`
		Order Order  `json:"o"`
		Value any    `json:"v"`
	}
)

func (t PageToken) Columns() []Column { return t.cols }

// Encrypt encrypts the page token using the first key in the provided keyset.
// It panics if no keys are provided.
func (t PageToken) Encrypt(keys [][32]byte) string {
	if len(keys) == 0 {
		panic("keyset pagination: cannot encrypt page token with no keys")
	}
	return hyrumtoken.Marshal(&keys[0], t)
}

func (t PageToken) MarshalJSON() ([]byte, error) {
	now := time.Now
	if t.testNow != nil {
		now = t.testNow
	}
	toEncode := jsonPageToken{
		ExpiresAt: now().Add(time.Hour).UTC(),
		Cols:      t.cols,
	}
	return json.Marshal(toEncode)
}

var ErrPageTokenExpired = herodot.ErrBadRequest.WithReason("page token expired, do not persist page tokens")

func (t *PageToken) UnmarshalJSON(data []byte) error {
	rawToken := jsonPageToken{}
	if err := json.Unmarshal(data, &rawToken); err != nil {
		return err
	}
	t.cols = rawToken.Cols
	now := time.Now
	if t.testNow != nil {
		now = t.testNow
	}
	if rawToken.ExpiresAt.Before(now().UTC()) {
		return errors.WithStack(ErrPageTokenExpired)
	}
	return nil
}

func NewPageToken(cols ...Column) PageToken {
	return PageToken{cols: cols}
}
