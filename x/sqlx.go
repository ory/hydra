package x

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

type JoseJSONWebKeySet struct {
	// swagger:ignore
	*jose.JSONWebKeySet
}

func (n *JoseJSONWebKeySet) Scan(value interface{}) error {
	v := fmt.Sprintf("%s", value)
	if len(v) == 0 {
		return nil
	}
	return errors.WithStack(json.Unmarshal([]byte(v), n))
}

func (n *JoseJSONWebKeySet) Value() (driver.Value, error) {
	value, err := json.Marshal(n)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(value), nil
}

type JSONRawMessage json.RawMessage

func (m *JSONRawMessage) Scan(value interface{}) error {
	*m = []byte(fmt.Sprintf("%s", value))
	return nil
}

func (m JSONRawMessage) Value() (driver.Value, error) {
	return string(m), nil
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
