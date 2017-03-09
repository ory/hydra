package fosite

import (
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (f *Fosite) NewAccessResponse(ctx context.Context, req *http.Request, requester AccessRequester) (AccessResponder, error) {
	var err error
	var tk TokenEndpointHandler

	response := NewAccessResponse()
	for _, tk = range f.TokenEndpointHandlers {
		if err = tk.PopulateTokenEndpointResponse(ctx, req, requester, response); errors.Cause(err) == ErrUnknownRequest {
		} else if err != nil {
			return nil, err
		}
	}

	if response.GetAccessToken() == "" || response.GetTokenType() == "" {
		return nil, errors.Wrap(ErrServerError, "Access token or token type not set")
	}

	return response, nil
}
