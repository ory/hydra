// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlxx

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"github.com/pkg/errors"
)

// Duration represents a JSON and SQL compatible time.Duration.
// swagger:type string
type Duration time.Duration

// MarshalJSON returns m as the JSON encoding of m.
func (ns Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(ns).String())
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *Duration) UnmarshalJSON(data []byte) error {
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

// StringSliceJSONFormat represents []string{} which is encoded to/from JSON for SQL storage.
type StringSliceJSONFormat []string

// Scan implements the Scanner interface.
func (m *StringSliceJSONFormat) Scan(value interface{}) error {
	val := fmt.Sprintf("%s", value)
	if len(val) == 0 {
		val = "[]"
	}

	if parsed := gjson.Parse(val); parsed.Type == gjson.Null {
		val = "[]"
	} else if !parsed.IsArray() {
		return errors.Errorf("expected JSON value to be an array but got type: %s", parsed.Type.String())
	}

	return errors.WithStack(json.Unmarshal([]byte(val), &m))
}

// Value implements the driver Valuer interface.
func (m StringSliceJSONFormat) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}

	encoded, err := json.Marshal(&m)
	return string(encoded), errors.WithStack(err)
}

// StringSlicePipeDelimiter de/encodes the string slice to/from a SQL string.
type StringSlicePipeDelimiter []string

// Scan implements the Scanner interface.
func (n *StringSlicePipeDelimiter) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	*n = scanStringSlice('|', s.String)
	return nil
}

// Value implements the driver Valuer interface.
func (n StringSlicePipeDelimiter) Value() (driver.Value, error) {
	return valueStringSlice('|', n), nil
}

func scanStringSlice(delimiter rune, value interface{}) []string {
	escaped := false
	s := fmt.Sprintf("%s", value)
	splitted := strings.FieldsFunc(s, func(r rune) bool {
		if r == '\\' {
			escaped = !escaped
		} else if escaped && r != delimiter {
			escaped = false
		}
		return !escaped && r == delimiter
	})
	for k, v := range splitted {
		splitted[k] = strings.ReplaceAll(v, "\\"+string(delimiter), string(delimiter))
	}
	return splitted
}

func valueStringSlice(delimiter rune, value []string) string {
	replace := make([]string, len(value))
	for k, v := range value {
		replace[k] = strings.ReplaceAll(v, string(delimiter), "\\"+string(delimiter))
	}
	return strings.Join(replace, string(delimiter))
}

// NullBool represents a bool that may be null.
// NullBool implements the Scanner interface so
// swagger:type bool
// swagger:model nullBool
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBool) Scan(value interface{}) error {
	d := sql.NullBool{}
	if err := d.Scan(value); err != nil {
		return err
	}

	ns.Bool = d.Bool
	ns.Valid = d.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (ns NullBool) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Bool)
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *NullBool) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	ns.Valid = true
	return errors.WithStack(json.Unmarshal(data, &ns.Bool))
}

// FalsyNullBool represents a bool that may be null.
// It JSON decodes to false if null.
//
// swagger:type bool
// swagger:model falsyNullBool
type FalsyNullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (ns *FalsyNullBool) Scan(value interface{}) error {
	d := sql.NullBool{}
	if err := d.Scan(value); err != nil {
		return err
	}

	ns.Bool = d.Bool
	ns.Valid = d.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (ns FalsyNullBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (ns FalsyNullBool) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("false"), nil
	}
	return json.Marshal(ns.Bool)
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *FalsyNullBool) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	ns.Valid = true
	return errors.WithStack(json.Unmarshal(data, &ns.Bool))
}

// swagger:type string
// swagger:model nullString
type NullString string

// MarshalJSON returns m as the JSON encoding of m.
func (ns NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ns))
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	if len(data) == 0 {
		return nil
	}
	return errors.WithStack(json.Unmarshal(data, (*string)(ns)))
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	var v sql.NullString
	if err := (&v).Scan(value); err != nil {
		return err
	}
	*ns = NullString(v.String)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if len(ns) == 0 {
		return sql.NullString{}.Value()
	}
	return sql.NullString{Valid: true, String: string(ns)}.Value()
}

// String implements the Stringer interface.
func (ns NullString) String() string {
	return string(ns)
}

// NullTime implements sql.NullTime functionality.
//
// swagger:model nullTime
// required: false
type NullTime time.Time

// Scan implements the Scanner interface.
func (ns *NullTime) Scan(value interface{}) error {
	var v sql.NullTime
	if err := (&v).Scan(value); err != nil {
		return err
	}
	*ns = NullTime(v.Time)
	return nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (ns NullTime) MarshalJSON() ([]byte, error) {
	var t *time.Time
	if !time.Time(ns).IsZero() {
		tt := time.Time(ns)
		t = &tt
	}
	return json.Marshal(t)
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *NullTime) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*ns = NullTime(t)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	return sql.NullTime{Valid: !time.Time(ns).IsZero(), Time: time.Time(ns)}.Value()
}

// MapStringInterface represents a map[string]interface that works well with JSON, SQL, and Swagger.
type MapStringInterface map[string]interface{}

// Scan implements the Scanner interface.
func (n *MapStringInterface) Scan(value interface{}) error {
	v := fmt.Sprintf("%s", value)
	if len(v) == 0 {
		return nil
	}
	return errors.WithStack(json.Unmarshal([]byte(v), n))
}

// Value implements the driver Valuer interface.
func (n MapStringInterface) Value() (driver.Value, error) {
	value, err := json.Marshal(n)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(value), nil
}

// JSONArrayRawMessage represents a json.RawMessage which only accepts arrays that works well with JSON, SQL, and Swagger.
type JSONArrayRawMessage json.RawMessage

// Scan implements the Scanner interface.
func (m *JSONArrayRawMessage) Scan(value interface{}) error {
	val := fmt.Sprintf("%s", value)
	if len(val) == 0 {
		val = "[]"
	}

	if parsed := gjson.Parse(val); parsed.Type == gjson.Null {
		val = "[]"
	} else if !parsed.IsArray() {
		return errors.Errorf("expected JSON value to be an array but got type: %s", parsed.Type.String())
	}

	*m = []byte(val)
	return nil
}

// Value implements the driver Valuer interface.
func (m JSONArrayRawMessage) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}

	if parsed := gjson.ParseBytes(m); parsed.Type == gjson.Null {
		return "[]", nil
	} else if !parsed.IsArray() {
		return nil, errors.Errorf("expected JSON value to be an array but got type: %s", parsed.Type.String())
	}

	return string(m), nil
}

// JSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger.
type JSONRawMessage json.RawMessage

// Scan implements the Scanner interface.
func (m *JSONRawMessage) Scan(value interface{}) error {
	*m = []byte(fmt.Sprintf("%s", value))
	return nil
}

// Value implements the driver Valuer interface.
func (m JSONRawMessage) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "null", nil
	}
	return string(m), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m JSONRawMessage) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
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

// NullJSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger and is NULLable-
//
// swagger:model nullJsonRawMessage
type NullJSONRawMessage json.RawMessage

// Scan implements the Scanner interface.
func (m *NullJSONRawMessage) Scan(value interface{}) error {
	if value == nil {
		value = "null"
	}
	*m = []byte(fmt.Sprintf("%s", value))
	return nil
}

// Value implements the driver Valuer interface.
func (m NullJSONRawMessage) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return string(m), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m NullJSONRawMessage) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *NullJSONRawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// JSONScan is a generic helper for storing a value as a JSON blob in SQL.
func JSONScan(dst interface{}, value interface{}) error {
	if value == nil {
		value = "null"
	}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", value)), &dst); err != nil {
		return fmt.Errorf("unable to decode payload to: %s", err)
	}
	return nil
}

// JSONValue is a generic helper for retrieving a SQL JSON-encoded value.
func JSONValue(src interface{}) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&src); err != nil {
		return nil, err
	}
	return b.String(), nil
}

// NullInt64 represents an int64 that may be null.
// swagger:model nullInt64
type NullInt64 struct {
	Int   int64
	Valid bool // Valid is true if Duration is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullInt64) Scan(value interface{}) error {
	d := sql.NullInt64{}
	if err := d.Scan(value); err != nil {
		return err
	}

	ns.Int = d.Int64
	ns.Valid = d.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullInt64) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Int, nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (ns NullInt64) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Int)
}

// UnmarshalJSON sets *m to a copy of data.
func (ns *NullInt64) UnmarshalJSON(data []byte) error {
	if ns == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	ns.Valid = true
	return errors.WithStack(json.Unmarshal(data, &ns.Int))
}

// NullDuration represents a nullable JSON and SQL compatible time.Duration.
//
// swagger:type string
// swagger:model nullDuration
type NullDuration struct {
	Duration time.Duration
	Valid    bool
}

// Scan implements the Scanner interface.
func (ns *NullDuration) Scan(value interface{}) error {
	d := sql.NullInt64{}
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

func (ns Duration) IsZero() bool                { return ns == 0 }
func (m StringSliceJSONFormat) IsZero() bool    { return len(m) == 0 }
func (n StringSlicePipeDelimiter) IsZero() bool { return len(n) == 0 }
func (ns NullBool) IsZero() bool                { return !ns.Valid || !ns.Bool }
func (ns FalsyNullBool) IsZero() bool           { return !ns.Valid || !ns.Bool }
func (ns NullString) IsZero() bool              { return len(ns) == 0 }
func (ns NullTime) IsZero() bool                { return time.Time(ns).IsZero() }
func (n MapStringInterface) IsZero() bool       { return len(n) == 0 }
func (m JSONArrayRawMessage) IsZero() bool      { return len(m) == 0 || string(m) == "[]" }
func (m JSONRawMessage) IsZero() bool           { return len(m) == 0 || string(m) == "null" }
func (m NullJSONRawMessage) IsZero() bool       { return len(m) == 0 || string(m) == "null" }
func (ns NullInt64) IsZero() bool               { return !ns.Valid || ns.Int == 0 }
func (ns NullDuration) IsZero() bool            { return !ns.Valid || ns.Duration == 0 }
