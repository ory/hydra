// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/errorsx"

	jose "github.com/go-jose/go-jose/v3"
)

// swagger:type JSONWebKeySet
type JoseJSONWebKeySet struct {
	// swagger:ignore
	*jose.JSONWebKeySet
}

func (n *JoseJSONWebKeySet) Scan(value interface{}) error {
	v := fmt.Sprintf("%s", value)
	if len(v) == 0 {
		return nil
	}
	return errorsx.WithStack(json.Unmarshal([]byte(v), n))
}

func (n *JoseJSONWebKeySet) Value() (driver.Value, error) {
	value, err := json.Marshal(n)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}
	return string(value), nil
}

type Duration time.Duration

// MarshalJSON returns m as the JSON encoding of m.
func (ns Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(ns).String())
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *Duration) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	p, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*ns = Duration(p)
	return nil
}

// swagger:model NullDuration
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type swaggerNullDuration string

// NullDuration represents a nullable JSON and SQL compatible time.Duration.
//
// TODO delete this type and replace it with ory/x/sqlxx/NullDuration when applying the custom client token TTL patch to Hydra 2.x
//
// swagger:ignore
type NullDuration struct {
	Duration time.Duration
	Valid    bool
}

// Scan implements the Scanner interface.
func (ns *NullDuration) Scan(value interface{}) error {
	var d = sql.NullInt64{}
	if err := d.Scan(value); err != nil {
		return err
	}

	ns.Duration = time.Duration(d.Int64)
	ns.Valid = d.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullDuration) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return int64(ns.Duration), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (ns NullDuration) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ns.Duration.String())
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *NullDuration) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	p, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	ns.Duration = p
	ns.Valid = true
	return nil
}
