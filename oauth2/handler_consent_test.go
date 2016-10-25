package oauth2

import (
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerConsent(t *testing.T) {
	h := new(Handler)
	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + "/oauth2/consent")
	assert.Nil(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.NotEmpty(t, body)
}
