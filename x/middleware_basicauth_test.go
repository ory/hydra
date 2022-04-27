package x

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicAuthMiddleware_ServeHTTP(t *testing.T) {
	tests := []struct {
		name             string
		requiredUsername string
		requiredPassword string
		givenUsername    string
		givenPassword    string
		givenNoAuth      bool
		expectedNext     bool
		expectedStatus   int
	}{{
		name:             "valid username/password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenUsername:    "foo",
		givenPassword:    "bar",
		expectedNext:     true,
		expectedStatus:   http.StatusOK,
	}, {
		name:             "valid empty username",
		requiredUsername: "",
		requiredPassword: "bar",
		givenUsername:    "",
		givenPassword:    "bar",
		expectedNext:     true,
		expectedStatus:   http.StatusOK,
	}, {
		name:             "valid empty password",
		requiredUsername: "foo",
		requiredPassword: "",
		givenUsername:    "foo",
		givenPassword:    "",
		expectedNext:     true,
		expectedStatus:   http.StatusOK,
	}, {
		name:             "valid empty username/password",
		requiredUsername: "",
		requiredPassword: "",
		givenUsername:    "",
		givenPassword:    "",
		expectedNext:     true,
		expectedStatus:   http.StatusOK,
	}, {
		name:             "no username/password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "invalid username/password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenUsername:    "foo",
		givenPassword:    "not-bar",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "invalid username, valid password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenUsername:    "not-foo",
		givenPassword:    "bar",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "valid username, invalid password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenUsername:    "foo",
		givenPassword:    "not-bar",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "invalid empty username",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenUsername:    "",
		givenPassword:    "bar",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "invalid empty password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenUsername:    "foo",
		givenPassword:    "",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "invalid empty username/password",
		requiredUsername: "",
		requiredPassword: "",
		givenUsername:    "foo",
		givenPassword:    "bar",
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "missing username/password",
		requiredUsername: "foo",
		requiredPassword: "bar",
		givenNoAuth:      true,
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}, {
		name:             "missing auth with empty required username/password",
		requiredUsername: "",
		requiredPassword: "",
		givenNoAuth:      true,
		expectedNext:     false,
		expectedStatus:   http.StatusUnauthorized,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewBasicAuthMiddleware(tt.requiredUsername, tt.requiredPassword)

			// add a next handler that sets nextCalled to true
			nextCalled := false
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				nextCalled = true
			})

			// construct a request
			req, err := http.NewRequest("GET", "", nil)
			require.NoError(t, err)

			// set the username and password
			if !tt.givenNoAuth {
				req.SetBasicAuth(tt.givenUsername, tt.givenPassword)
			}

			// serve the request
			rw := httptest.NewRecorder()
			m.ServeHTTP(rw, req, next.ServeHTTP)

			// check the response and whether the next handler was called
			require.Equal(t, tt.expectedNext, nextCalled)
			require.Equal(t, tt.expectedStatus, rw.Code)
		})
	}
}
