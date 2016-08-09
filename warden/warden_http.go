package warden

import (
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type HTTPWarden struct {
	Client   *http.Client
	Dry      bool
	Endpoint *url.URL
}

func (w *HTTPWarden) TokenFromRequest(r *http.Request) string {
	return fosite.AccessTokenFromRequest(r)
}

func (w *HTTPWarden) SetClient(c *clientcredentials.Config) {
	w.Client = c.Client(oauth2.NoContext)
}

func (w *HTTPWarden) IntrospectToken(ctx context.Context, token string) (*firewall.Introspection, error) {
	var resp = new(firewall.Introspection)
	var ep = *w.Endpoint
	ep.Path = IntrospectPath
	agent := &pkg.SuperAgent{URL: ep.String(), Client: w.Client}

	data := url.Values{"token": []string{token}}
	hreq, err := http.NewRequest("POST", ep.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, errors.New(err)
	}

	hreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hreq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	hres, err := w.Client.Do(hreq)
	if err != nil {
		return nil, errors.New(err)
	}

	if hres.StatusCode < 200 || hres.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(hres.Body)
		return nil, errors.Errorf("Expected 2xx status code but got %d.\n%s", hres.StatusCode, body)
	} else if err := json.NewDecoder(hres.Body).Decode(resp); err != nil {
		body, _ := ioutil.ReadAll(hres.Body)
		return nil, errors.Errorf("%s: %s", err, body)
	}

	if err := agent.POST(&struct {
		Token string `json:"token"`
	}{Token: token}, &hres); err != nil {
		return nil, err
	} else if !resp.Active {
		return nil, errors.New("Token is malformed, expired or otherwise invalid")
	}

	return resp, nil
}

func (w *HTTPWarden) TokenAllowed(ctx context.Context, token string, a *ladon.Request, scopes ...string) (*firewall.Context, error) {
	var resp = struct {
		*firewall.Context
		Allowed bool `json:"allowed"`
	}{}

	var ep = *w.Endpoint
	ep.Path = TokenAllowedHandlerPath
	agent := &pkg.SuperAgent{URL: ep.String(), Client: w.Client}
	if err := agent.POST(&WardenAccessRequest{
		WardenAuthorizedRequest: &WardenAuthorizedRequest{
			Token:  token,
			Scopes: scopes,
		},
		Request: a,
	}, &resp); err != nil {
		return nil, err
	} else if !resp.Allowed {
		return nil, errors.New("Token is not valid")
	}

	return resp.Context, nil
}

func (w *HTTPWarden) IsAllowed(ctx context.Context, a *ladon.Request) error {
	var allowed = struct {
		Allowed bool `json:"allowed"`
	}{}

	var ep = *w.Endpoint
	ep.Path = AllowedHandlerPath
	agent := &pkg.SuperAgent{URL: ep.String(), Client: w.Client}
	if err := agent.POST(a, &allowed); err != nil {
		return err
	} else if !allowed.Allowed {
		return errors.New("Forbidden")
	}

	return nil
}

func (w *HTTPWarden) InspectToken(ctx context.Context, token string, scopes ...string) (*firewall.Context, error) {
	var resp = struct {
		*firewall.Context
		Valid bool `json:"valid"`
	}{}

	var ep = *w.Endpoint
	ep.Path = TokenValidHandlerPath
	agent := &pkg.SuperAgent{URL: ep.String(), Client: w.Client}
	if err := agent.POST(&WardenAuthorizedRequest{
		Token:  token,
		Scopes: scopes,
	}, &resp); err != nil {
		return nil, err
	} else if !resp.Valid {
		return nil, errors.New("Token is not valid")
	}

	return resp.Context, nil
}
