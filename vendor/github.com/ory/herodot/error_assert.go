package herodot

import "github.com/pkg/errors"

func toError(e interface{}) error {
	if e == nil {
		return errors.New("Error passed to WriteErrorCode is nil")
	}

	err, ok := e.(error)
	if !ok {
		return errors.New("Error passed to WriteErrorCode does not implement the error interface")
	}

	return err
}
