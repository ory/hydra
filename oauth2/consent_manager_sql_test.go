package oauth2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsentRequestSqlDataTransforms(t *testing.T) {
	for _, tc := range []struct {
		d string
		r *ConsentRequest
	}{
		{
			d: "fully hydrated request object",
			r: &ConsentRequest{
				ID:               "id",
				ClientID:         "client-id",
				RequestedScopes:  []string{"foo", "bar"},
				GrantedScopes:    []string{"baz", "bar"},
				CSRF:             "some-csrf",
				ExpiresAt:        time.Now().Round(time.Second),
				Consent:          ConsentRequestAccepted,
				DenyReason:       "some reason",
				AccessTokenExtra: map[string]interface{}{"atfoo": "bar", "atbaz": "bar"},
				IDTokenExtra:     map[string]interface{}{"idfoo": "bar", "idbaz": "bar"},
				RedirectURL:      "https://redirect-me/foo",
				Subject:          "Peter",
			},
		},
	} {
		t.Run(tc.d, func(t *testing.T) {
			s, err := newConsentRequestSqlData(tc.r)
			require.Nil(t, err)

			o, err := s.toConsentRequest()
			require.NoError(t, err)

			assert.EqualValues(t, tc.r, o)
		})
	}
}
