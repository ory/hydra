package x

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
	jose "gopkg.in/square/go-jose.v2"
)

type JSONWebKeySet struct {
	*jose.JSONWebKeySet
}

func (n *JSONWebKeySet) Scan(value interface{}) error {
	v := fmt.Sprintf("%s", value)
	if len(v) == 0 {
		return nil
	}
	return errors.WithStack(json.Unmarshal([]byte(v), n))
}

func (n JSONWebKeySet) Value() (driver.Value, error) {
	value, err := json.Marshal(&n)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(value), nil
}

type JSONRawMessage json.RawMessage

func (n *JSONRawMessage) Scan(value interface{}) error {
	*n = []byte(fmt.Sprintf("%s",value))
	return nil
}

func (n JSONRawMessage) Value() (driver.Value, error) {
	return string(n), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m JSONRawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JSONRawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
