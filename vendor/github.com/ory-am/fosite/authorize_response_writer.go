package fosite

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (o *Fosite) NewAuthorizeResponse(ctx context.Context, r *http.Request, ar AuthorizeRequester, session Session) (AuthorizeResponder, error) {
	var resp = &AuthorizeResponse{
		Header:   http.Header{},
		Query:    url.Values{},
		Fragment: url.Values{},
	}

	ar.SetSession(session)
	for _, h := range o.AuthorizeEndpointHandlers {
		if err := h.HandleAuthorizeEndpointRequest(ctx, r, ar, resp); err != nil {
			return nil, err
		}
	}

	if !ar.DidHandleAllResponseTypes() {
		return nil, errors.WithStack(ErrUnsupportedResponseType)
	}

	return resp, nil
}
