package x

import (
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/pkg/errors"
	"net/http"
)

type enhancedError struct {
	*fosite.RFC6749Error
	RequestID string `json:"request_id"`
}

func ErrorEnhancer(r *http.Request, err error) interface{} {
	if e := new(herodot.DefaultError); errors.As(err, &e) {
		return &enhancedError{
			RFC6749Error: (&fosite.RFC6749Error{
				ErrorField:       e.Error(),
				DescriptionField: e.Reason(),
				CodeField:        e.StatusCode(),
			}).WithTrace(err),
			RequestID: r.Header.Get("X-Request-Id"),
		}
	}

	return &enhancedError{
		RFC6749Error: fosite.ErrorToRFC6749Error(err),
		RequestID: r.Header.Get("X-Request-Id"),
	}
}
