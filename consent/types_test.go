package consent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
)

func TestToRFCError(t *testing.T) {
	for k, tc := range []struct {
		input  *RequestDeniedError
		expect *fosite.RFC6749Error
	}{
		{
			input: &RequestDeniedError{
				Name:  "not empty",
				valid: true,
			},
			expect: &fosite.RFC6749Error{
				Name:        "not empty",
				Description: "",
				Code:        fosite.ErrInvalidRequest.Code,
				Debug:       "",
			},
		},
		{
			input: &RequestDeniedError{
				Name:        "",
				Description: "not empty",
				valid:       true,
			},
			expect: &fosite.RFC6749Error{
				Name:        "request was denied",
				Description: "not empty",
				Code:        fosite.ErrInvalidRequest.Code,
				Debug:       "",
			},
		},
		{
			input: &RequestDeniedError{valid: true},
			expect: &fosite.RFC6749Error{
				Name:        "request was denied",
				Description: "",
				Hint:        "",
				Code:        fosite.ErrInvalidRequest.Code,
				Debug:       "",
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			require.EqualValues(t, tc.input.toRFCError(), tc.expect)
		})
	}
}

func TestRequestDeniedError(t *testing.T) {
	var e *RequestDeniedError
	v, err := e.Value()
	require.NoError(t, err)
	assert.EqualValues(t, "{}", fmt.Sprintf("%v", v))
}
