// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/ssoready/hyrumtoken"

	"github.com/ory/herodot"
)

var fallbackEncryptionKey = &[32]byte{}

type (
	PageToken struct {
		testNow func() time.Time
		cols    []Column
	}
	jsonPageToken = struct {
		ExpiresAt time.Time    `json:"e"`
		Cols      []jsonColumn `json:"c"`
	}
	jsonColumn = struct {
		Name      string    `json:"n"`
		Order     Order     `json:"o"`
		ValueAny  any       `json:"v"`
		ValueTime time.Time `json:"vt"`
		ValueUUID uuid.UUID `json:"vu"`
		ValueInt  int64     `json:"vi"`
	}
	Column struct {
		Name  string
		Order Order
		Value any
	}
)

func (t PageToken) Columns() []Column { return t.cols }

// Encrypt encrypts the page token using the first key in the provided keyset.
// It uses a fallback key if no keys are provided.
func (t PageToken) Encrypt(keys [][32]byte) string {
	key := fallbackEncryptionKey
	if len(keys) > 0 {
		key = &keys[0]
	}
	return hyrumtoken.Marshal(key, t)
}

func (t PageToken) MarshalJSON() ([]byte, error) {
	now := time.Now
	if t.testNow != nil {
		now = t.testNow
	}
	toEncode := jsonPageToken{
		ExpiresAt: now().Add(time.Hour).UTC(),
		Cols:      make([]jsonColumn, len(t.cols)),
	}
	for i, col := range t.cols {
		toEncode.Cols[i] = jsonColumn{
			Name:  col.Name,
			Order: col.Order,
		}
		switch v := col.Value.(type) {
		case time.Time:
			toEncode.Cols[i].ValueTime = v
		case uuid.UUID:
			toEncode.Cols[i].ValueUUID = v
		case int64:
			toEncode.Cols[i].ValueInt = v
		default:
			toEncode.Cols[i].ValueAny = v
		}
	}
	return json.Marshal(toEncode)
}

var ErrPageTokenExpired = herodot.ErrBadRequest.WithReason("page token expired, do not persist page tokens")

func (t *PageToken) UnmarshalJSON(data []byte) error {
	rawToken := jsonPageToken{}
	if err := json.Unmarshal(data, &rawToken); err != nil {
		return err
	}
	t.cols = make([]Column, len(rawToken.Cols))
	for i, col := range rawToken.Cols {
		t.cols[i] = Column{
			Name:  col.Name,
			Order: col.Order,
		}
		switch {
		case col.ValueAny != nil:
			t.cols[i].Value = col.ValueAny
		case !col.ValueTime.IsZero():
			t.cols[i].Value = col.ValueTime
		case col.ValueUUID != uuid.Nil:
			t.cols[i].Value = col.ValueUUID
		case col.ValueInt != 0:
			t.cols[i].Value = col.ValueInt
		}
	}
	now := time.Now
	if t.testNow != nil {
		now = t.testNow
	}
	if rawToken.ExpiresAt.Before(now().UTC()) {
		return errors.WithStack(ErrPageTokenExpired)
	}
	return nil
}

func NewPageToken(cols ...Column) PageToken { return PageToken{cols: cols} }
