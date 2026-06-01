// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/v2"

	"github.com/pkg/errors"

	stdjson "encoding/json"
)

// KoanfMemory implements a KoanfMemory provider.
type KoanfMemory struct {
	doc    stdjson.RawMessage
	parser koanf.Parser
}

// NewKoanfMemory returns a file provider.
func NewKoanfMemory(doc stdjson.RawMessage) *KoanfMemory {
	return &KoanfMemory{
		doc:    doc,
		parser: json.Parser(),
	}
}

func (f *KoanfMemory) SetDoc(doc stdjson.RawMessage) {
	f.doc = doc
}

// ReadBytes reads the contents of a file on disk and returns the bytes.
func (f *KoanfMemory) ReadBytes() ([]byte, error) {
	return nil, errors.New("file provider does not support this method")
}

// Read is not supported by the file provider.
func (f *KoanfMemory) Read() (map[string]interface{}, error) {
	v, err := f.parser.Unmarshal(f.doc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return v, nil
}
