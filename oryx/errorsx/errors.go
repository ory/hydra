// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package errorsx

import (
	"github.com/pkg/errors"

	"github.com/ory/herodot"
)

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
// Deprecated: you should probably use errors.As instead.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok || cause.Cause() == nil {
			break
		}
		err = cause.Cause()
	}
	return err
}

// WithStack mirror pkg/errors.WithStack but does not wrap existing stack
// traces.
// Deprecated: you should probably use errors.WithStack instead and only annotate stacks when it makes sense.
func WithStack(err error) error {
	if e, ok := err.(StackTracer); ok && len(e.StackTrace()) > 0 {
		return err
	}

	return errors.WithStack(err)
}

// StatusCodeCarrier can be implemented by an error to support setting status codes in the error itself.
type StatusCodeCarrier interface {
	// StatusCode returns the status code of this error.
	StatusCode() int
}

// RequestIDCarrier can be implemented by an error to support error contexts.
type RequestIDCarrier interface {
	// RequestID returns the ID of the request that caused the error, if applicable.
	RequestID() string
}

// ReasonCarrier can be implemented by an error to support error contexts.
type ReasonCarrier interface {
	// Reason returns the reason for the error, if applicable.
	Reason() string
}

// DebugCarrier can be implemented by an error to support error contexts.
type DebugCarrier interface {
	// Debug returns debugging information for the error, if applicable.
	Debug() string
}

// StatusCarrier can be implemented by an error to support error contexts.
type StatusCarrier interface {
	// ID returns the error id, if applicable.
	Status() string
}

// DetailsCarrier can be implemented by an error to support error contexts.
type DetailsCarrier interface {
	// Details returns details on the error, if applicable.
	Details() map[string]interface{}
}

// IDCarrier can be implemented by an error to support error contexts.
type IDCarrier interface {
	// ID returns application error ID on the error, if applicable.
	ID() string
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func GetCodeFromHerodotError(err error) (code int, ok bool) {
	herodotErr := &herodot.DefaultError{}
	isHerodot := errors.As(err, &herodotErr)

	return herodotErr.CodeField, isHerodot
}
