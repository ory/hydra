package consent

import (
	"fmt"
	"testing"

	"github.com/ory/fosite"
	"github.com/stretchr/testify/require"
)

func TestToRFCError(t *testing.T) {
	for k, tc := range []struct {
		input  *RequestDeniedError
		expect *fosite.RFC6749Error
	}{
		{
			input: &RequestDeniedError{
				Description: "not empty",
			},
			expect: &fosite.RFC6749Error{
				Name:        "",
				Description: "not empty",
				Hint:        "",
				Code:        fosite.ErrInvalidRequest.Code,
				Debug:       "",
			},
		},
		{
			input: &RequestDeniedError{
				Name:        "not empty",
				Description: "not empty",
			},
			expect: &fosite.RFC6749Error{
				Name:        "not empty",
				Description: "not empty",
				Hint:        "",
				Code:        fosite.ErrInvalidRequest.Code,
				Debug:       "",
			},
		},
		{
			input: &RequestDeniedError{
				Description: "not empty",
				Hint:        "not empty",
			},
			expect: &fosite.RFC6749Error{
				Name:        "",
				Description: "not empty",
				Hint:        "not empty",
				Code:        fosite.ErrInvalidRequest.Code,
				Debug:       "",
			},
		},
		{
			input: &RequestDeniedError{
				Name: "not empty",
			},
			expect: &fosite.RFC6749Error{
				Name:  "not empty",
				Code:  fosite.ErrInvalidRequest.Code,
				Debug: "",
			},
		},
		{
			input: &RequestDeniedError{},
			expect: &fosite.RFC6749Error{
				Name:        requestDeniedErrorName,
				Description: requestDeniedErrorDescription,
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
