// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonschemax

import (
	"github.com/ory/jsonschema/v3"
)

// ErrorType is the schema error type.
type ErrorType int

const (
	// ErrorTypeMissing represents a validation that failed because a value is missing.
	ErrorTypeMissing ErrorType = iota + 1
)

// Error represents a schema error.
type Error struct {
	// Type is the error type.
	Type ErrorType

	// DocumentPointer is the JSON Pointer in the document.
	DocumentPointer string

	// SchemaPointer is the JSON Pointer in the schema.
	SchemaPointer string

	// DocumentFieldName is a pointer to the document in dot-notation: fo.bar.baz
	DocumentFieldName string
}

// NewFromSanthoshError converts github.com/santhosh-tekuri/jsonschema.ValidationError to Error.
func NewFromSanthoshError(validationError jsonschema.ValidationError) *Error {
	return &Error{
		// DocumentPointer:   JSONPointerToDotNotation(validationError.InstancePtr),
		// SchemaPointer:     JSONPointerToDotNotation(validationError.SchemaPtr),
		// DocumentFieldName: JSONPointerToDotNotation(validationError.InstancePtr),
	}
}
