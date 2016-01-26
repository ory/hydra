package signin

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mock = &signin{
	id:         "123",
	login:      "http://localhost:3000/login",
	redirectTo: "http://foobar.com",
}

func TestNew(t *testing.T) {
	m := New("321", mock.login, "http://foobar.com")
	assert.Equal(t, "321", m.id)
	assert.Equal(t, "http://localhost:3000/login", m.login)
	assert.Equal(t, "http://foobar.com", m.redirectTo)
}

func TestGetAuthCodeURL(t *testing.T) {
	require.NotEmpty(t, mock.GetAuthenticationURL("state"))
}

func TestExchangeCode(t *testing.T) {
	router := mux.NewRouter()
	user := "foouser"
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fmt.Sprintf(`{"subject": "%s"}`, r.URL.Query().Get("verify")))
	})
	ts := httptest.NewServer(router)
	mock.login = ts.URL + "/login"
	ses, err := mock.FetchSession(user)
	require.Nil(t, err)
	assert.Equal(t, user, ses.GetForcedLocalSubject())
}
