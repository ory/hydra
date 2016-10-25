package oauth2

import (
	"bytes"
	"encoding/json"
	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HTTPIntrospector struct {
	Client   *http.Client
	Dry      bool
	Endpoint *url.URL
}

func (i *HTTPIntrospector) TokenFromRequest(r *http.Request) string {
	return fosite.AccessTokenFromRequest(r)
}

func (i *HTTPIntrospector) SetClient(c *clientcredentials.Config) {
	i.Client = c.Client(oauth2.NoContext)
}

// IntrospectToken is capable of introspecting tokens according to https://tools.ietf.org/html/rfc7662
//
// The HTTP API is documented at http://docs.hdyra.apiary.io/#reference/oauth2/oauth2-token-introspection
func (i *HTTPIntrospector) IntrospectToken(ctx context.Context, token string, scopes ...string) (*Introspection, error) {
	var resp = &Introspection{
		Extra: make(map[string]interface{}),
	}
	var ep = *i.Endpoint
	ep.Path = IntrospectPath

	data := url.Values{"token": []string{token}, "scope": []string{strings.Join(scopes, " ")}}
	hreq, err := http.NewRequest("POST", ep.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	hreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hreq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	hres, err := i.Client.Do(hreq)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer hres.Body.Close()

	body, _ := ioutil.ReadAll(hres.Body)
	if hres.StatusCode < 200 || hres.StatusCode >= 300 {
		return nil, errors.Errorf("Expected 2xx status code but got %d.\n%s", hres.StatusCode, body)
	} else if err := json.Unmarshal(body, resp); err != nil {
		return nil, errors.Errorf("%s: %s", err, body)
	} else if !resp.Active {
		return nil, errors.New("Token is malformed, expired or otherwise invalid")
	}
	return resp, nil
}
