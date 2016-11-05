package warden

import (
	"net/http"
	"net/url"

	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
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

// TokenAllowed checks if a token is valid and if the token owner is allowed to perform an action on a resource.
// This endpoint requires a token, a scope, a resource name, an action name and a context.
//
// The HTTP API is documented at http://docs.hdyra.apiary.io/#reference/warden:-access-control-for-resource-providers/check-if-an-access-tokens-subject-is-allowed-to-do-something
func (w *HTTPWarden) TokenAllowed(ctx context.Context, token string, a *firewall.TokenAccessRequest, scopes ...string) (*firewall.Context, error) {
	var resp = struct {
		*firewall.Context
		Allowed bool `json:"allowed"`
	}{}

	var ep = *w.Endpoint
	ep.Path = TokenAllowedHandlerPath
	agent := &pkg.SuperAgent{URL: ep.String(), Client: w.Client}
	if err := agent.POST(&wardenAccessRequest{
		wardenAuthorizedRequest: &wardenAuthorizedRequest{
			Token:  token,
			Scopes: scopes,
		},
		TokenAccessRequest: a,
	}, &resp); err != nil {
		return nil, err
	} else if !resp.Allowed {
		return nil, errors.New("Token is not valid")
	}

	return resp.Context, nil
}

// IsAllowed checks if an arbitrary subject is allowed to perform an action on a resource.
//
// The HTTP API is documented at http://docs.hdyra.apiary.io/#reference/warden:-access-control-for-resource-providers/check-if-a-subject-is-allowed-to-do-something
func (w *HTTPWarden) IsAllowed(ctx context.Context, a *firewall.AccessRequest) error {
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
