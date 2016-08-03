package warden

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	. "github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/pkg/helper"
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

func (w *HTTPWarden) SetClient(c *clientcredentials.Config) {
	w.Client = c.Client(oauth2.NoContext)
}

func (w *HTTPWarden) IntrospectToken() {

}

func (w *HTTPWarden) ActionAllowed(ctx context.Context, token string, a *ladon.Request, scopes ...string) (*Context, error) {
	return w.doRequest(TokenAllowedHandlerPath, &WardenAccessRequest{
		Request: a,
		WardenAuthorizedRequest: &WardenAuthorizedRequest{
			Token:  token,
			Scopes: scopes,
		},
	})
}

func (w *HTTPWarden) HTTPRequestAllowed(ctx context.Context, r *http.Request, a *ladon.Request, scopes ...string) (*Context, error) {
	token := TokenFromRequest(r)
	if token == "" {
		return nil, errors.New(pkg.ErrUnauthorized)
	}

	return w.ActionAllowed(ctx, token, a, scopes...)
}

func (w *HTTPWarden) Authorized(ctx context.Context, token string, scopes ...string) (*Context, error) {
	return w.doRequest(TokenValidHandlerPath, &WardenAuthorizedRequest{
		Token:  token,
		Scopes: scopes,
	})
}

func (w *HTTPWarden) HTTPAuthorized(ctx context.Context, r *http.Request, scopes ...string) (*Context, error) {
	token := TokenFromRequest(r)
	if token == "" {
		return nil, errors.New(pkg.ErrUnauthorized)
	}

	return w.Authorized(ctx, token, scopes...)
}

func (w *HTTPWarden) doDry(req *http.Request) error {
	return helper.DoDryRequest(w.Dry, req)
}

func (w *HTTPWarden) doRequest(path string, request interface{}) (*Context, error) {
	out, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New(err)
	}

	var ep = new(url.URL)
	*ep = *w.Endpoint
	ep.Path = path
	req, err := http.NewRequest("POST", ep.String(), bytes.NewBuffer(out))
	if err != nil {
		return nil, errors.New(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if err := w.doDry(req); err != nil {
		return nil, err
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, errors.New(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(err)
		}

		return nil, errors.Errorf("Got error (%d): %s", resp.StatusCode, all)
	}

	var epResp = struct {
		*Context
		Valid   bool `json:"valid"`
		Allowed bool `json:"allowed"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&epResp); err != nil {
		return nil, errors.New(err)
	}

	if epResp.Valid || epResp.Allowed {
		return epResp.Context, nil
	}

	return nil, errors.Errorf("Token subject has insufficient rights or invalid token")
}
