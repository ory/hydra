// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package spec

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
)

//go:embed config.json
var ConfigValidationSchema []byte

var ConfigSchemaID string

func init() {
	ConfigSchemaID = gjson.GetBytes(ConfigValidationSchema, "$id").String()
	if ConfigSchemaID == "" {
		ConfigSchemaID = uuid.Must(uuid.NewV4()).String() + ".json"
	}
}

// AddConfigSchema should be used instead of the schema itself to auto-register the dependencies schemas.
func AddConfigSchema(compiler interface {
	AddResource(url string, r io.Reader) error
}) error {
	if err := otelx.AddConfigSchema(compiler); err != nil {
		return err
	}
	if err := logrusx.AddConfigSchema(compiler); err != nil {
		return err
	}

	return errors.WithStack(compiler.AddResource(ConfigSchemaID, bytes.NewReader(ConfigValidationSchema)))
}
