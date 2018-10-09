package herodot

import (
	"net/http"

	"github.com/pkg/errors"
)

type DefaultError struct {
	CodeField    int                      `json:"code,omitempty"`
	StatusField  string                   `json:"status,omitempty"`
	RIDField     string                   `json:"request,omitempty"`
	ReasonField  string                   `json:"reason,omitempty"`
	DetailsField map[string][]interface{} `json:"details,omitempty"`
	ErrorField   string                   `json:"message"`
}

func (e *DefaultError) Status() string {
	return e.StatusField
}

func (e *DefaultError) Error() string {
	return e.ErrorField
}

func (e *DefaultError) RequestID() string {
	return e.RIDField
}

func (e *DefaultError) Reason() string {
	return e.ReasonField
}

func (e *DefaultError) Details() map[string][]interface{} {
	return e.DetailsField
}

func (e *DefaultError) StatusCode() int {
	return e.CodeField
}

func (e *DefaultError) WithReason(reason string) *DefaultError {
	err := *e
	err.ReasonField = reason
	return &err
}

func (e *DefaultError) WithDetail(key string, message ...interface{}) *DefaultError {
	err := *e
	if err.DetailsField == nil {
		err.DetailsField = map[string][]interface{}{}
	}
	err.DetailsField[key] = append(err.DetailsField[key], message...)
	return &err
}

func toDefaultError(err error, rid string) *DefaultError {
	if e, ok := errors.Cause(err).(ErrorContextCarrier); ok {
		id := e.RequestID()
		if id == "" {
			id = rid
		}

		return &DefaultError{
			CodeField:    e.StatusCode(),
			ReasonField:  e.Reason(),
			RIDField:     rid,
			ErrorField:   err.Error(),
			DetailsField: e.Details(),
			StatusField:  e.Status(),
		}
	} else if e, ok := errors.Cause(err).(StatusCodeCarrier); ok {
		return &DefaultError{
			RIDField:     rid,
			CodeField:    e.StatusCode(),
			ErrorField:   err.Error(),
			DetailsField: map[string][]interface{}{},
		}
	}

	return &DefaultError{
		RIDField:     rid,
		ErrorField:   err.Error(),
		CodeField:    http.StatusInternalServerError,
		DetailsField: map[string][]interface{}{},
	}
}
