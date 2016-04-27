package warden

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type HTTPWarden struct {
	Client *http.Client

	Endpoint *url.URL
}

func (w *HTTPWarden) SetClient(c *clientcredentials.Config) {
	w.Client = c.Client(oauth2.NoContext)
}

func (w *HTTPWarden) ActionAllowed(ctx context.Context, token string, a *ladon.Request, scopes ...string) (*Context, error) {
	return w.doRequest(AllowedHandlerPath, &WardenAccessRequest{
		Request: a,
		WardenAuthorizedRequest: &WardenAuthorizedRequest{
			InspectToken: token,
			Scopes:       scopes,
		},
	})
}

func (w *HTTPWarden) HTTPActionAllowed(ctx context.Context, r *http.Request, a *ladon.Request, scopes ...string) (*Context, error) {
	token := TokenFromRequest(r)
	if token == "" {
		return nil, errors.New(pkg.ErrUnauthorized)
	}

	return w.ActionAllowed(ctx, token, a, scopes...)
}

func (w *HTTPWarden) Authorized(ctx context.Context, token string, scopes ...string) (*Context, error) {
	return w.doRequest(AuthorizedHandlerPath, &WardenAuthorizedRequest{
		InspectToken: token,
		Scopes:       scopes,
	})
}

func (w *HTTPWarden) HTTPAuthorized(ctx context.Context, r *http.Request, scopes ...string) (*Context, error) {
	token := TokenFromRequest(r)
	if token == "" {
		return nil, errors.New(pkg.ErrUnauthorized)
	}

	return w.Authorized(ctx, token, scopes...)
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

	var epResp WardenResponse
	if err := json.NewDecoder(resp.Body).Decode(&epResp); err != nil {
		return nil, errors.New(err)
	}

	return epResp.Context, nil
}
