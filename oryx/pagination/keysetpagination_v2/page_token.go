// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"

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
		Name          string `json:"n"`
		Order         Order  `json:"o,omitempty"`
		Nullable      bool   `json:"nl,omitempty"`
		HasConstraint bool   `json:"hc,omitempty"`

		ValueAny   any        `json:"v,omitempty"`
		ValueTime  *time.Time `json:"vt,omitempty"`
		ValueUUID  *uuid.UUID `json:"vu,omitempty"`
		ValueInt64 *int64     `json:"vi,omitempty"`
		ValueNull  bool       `json:"vn,omitempty"`
	}
	Column struct {
		Name     string
		Order    Order
		Value    any
		Nullable bool
		// HasConstraint marks if the column is already constrained by WHERE condition.
		HasConstraint bool
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
	enc, err := t.encrypt(key)
	if err != nil {
		// This should basically never happen, only if reading from the random source or marshaling the token fails.
		// In both cases, we have a bigger problem than just not being able to generate the page token.
		// Therefore, if we do get an error, there is no point in returning it to the client,
		// as we already have a working result set. With this string, the next page will return an error,
		// but that is better than breaking the current page.
		return "internal error: failed to generate page token"
	}
	return enc
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
			Name:          col.Name,
			Order:         col.Order,
			Nullable:      col.Nullable,
			HasConstraint: col.HasConstraint,
		}
		switch v := col.Value.(type) {
		case time.Time:
			toEncode.Cols[i].ValueTime = new(v)
		case uuid.UUID:
			toEncode.Cols[i].ValueUUID = new(v)
		case uuid.NullUUID:
			if v.Valid {
				toEncode.Cols[i].ValueUUID = new(v.UUID)
			} else {
				toEncode.Cols[i].ValueNull = true
			}
		case sql.NullString:
			if v.Valid {
				toEncode.Cols[i].ValueAny = new(v.String)
			} else {
				toEncode.Cols[i].ValueNull = true
			}
		case int64:
			toEncode.Cols[i].ValueInt64 = new(v)
		case sql.NullInt64:
			if v.Valid {
				toEncode.Cols[i].ValueInt64 = new(v.Int64)
			} else {
				toEncode.Cols[i].ValueNull = true
			}
		default:
			toEncode.Cols[i].ValueAny = v
		}
	}
	return json.Marshal(toEncode)
}

func ErrPageTokenExpired() *herodot.DefaultError {
	return herodot.ErrBadRequest().WithError("page token expired, do not persist page tokens")
}

func (t *PageToken) UnmarshalJSON(data []byte) error {
	rawToken := jsonPageToken{}
	if err := json.Unmarshal(data, &rawToken); err != nil {
		return err
	}
	t.cols = make([]Column, len(rawToken.Cols))
	for i, col := range rawToken.Cols {
		t.cols[i] = Column{
			Name:          col.Name,
			Order:         col.Order,
			Nullable:      col.Nullable,
			HasConstraint: col.HasConstraint,
		}
		switch {
		case col.ValueNull:
			t.cols[i].Value = nil
		case col.ValueAny != nil:
			t.cols[i].Value = col.ValueAny
		case col.ValueUUID != nil:
			t.cols[i].Value = *col.ValueUUID
		case col.ValueTime != nil:
			t.cols[i].Value = *col.ValueTime
		case col.ValueInt64 != nil:
			t.cols[i].Value = *col.ValueInt64
		}
	}

	now := time.Now
	if t.testNow != nil {
		now = t.testNow
	}
	if rawToken.ExpiresAt.Before(now().UTC()) {
		return errors.WithStack(ErrPageTokenExpired())
	}
	return nil
}

func NewPageToken(cols ...Column) PageToken { return PageToken{cols: cols} }

func (t *PageToken) encrypt(key *[32]byte) (string, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", errors.Wrap(err, "cannot seed nonce")
	}

	raw, err := json.Marshal(t)
	if err != nil {
		return "", errors.Wrap(err, "cannot marshal page token")
	}

	enc := secretbox.Seal(nonce[:], raw, &nonce, key)
	return base64.URLEncoding.EncodeToString(enc), nil
}

func (t *PageToken) decrypt(key *[32]byte, s string) error {
	if s == "" {
		return errors.WithStack(ErrInvalidPaginationToken())
	}

	raw, err := base64.URLEncoding.DecodeString(s)
	if err != nil || len(raw) < 24 {
		return errors.WithStack(ErrInvalidPaginationToken())
	}

	var nonce [24]byte
	copy(nonce[:], raw[:24])

	dec, ok := secretbox.Open(nil, raw[24:], &nonce, key)
	if !ok {
		return errors.WithStack(ErrInvalidPaginationToken())
	}

	if err := json.Unmarshal(dec, t); err != nil {
		if errors.As(err, new(*herodot.DefaultError)) {
			return err
		}
		return errors.WithStack(herodot.ErrInternalServerError().WithReason("unable to unmarshal page token").WithDebug(err.Error()))
	}

	return nil
}
