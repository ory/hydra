package fosite

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// NewRevocationRequest handles incoming token revocation requests and
// validates various parameters as specified in:
// https://tools.ietf.org/html/rfc7009#section-2.1
//
// The authorization server first validates the client credentials (in
// case of a confidential client) and then verifies whether the token
// was issued to the client making the revocation request.  If this
// validation fails, the request is refused and the client is informed
// of the error by the authorization server as described below.
//
// In the next step, the authorization server invalidates the token.
// The invalidation takes place immediately, and the token cannot be
// used again after the revocation.
//
// * https://tools.ietf.org/html/rfc7009#section-2.2
// An invalid token type hint value is ignored by the authorization
// server and does not influence the revocation response.
func (f *Fosite) NewRevocationRequest(ctx context.Context, r *http.Request) error {
	if r.Method != "POST" {
		return errors.Wrap(ErrInvalidRequest, "HTTP method is not POST")
	} else if err := r.ParseForm(); err != nil {
		return errors.Wrap(ErrInvalidRequest, err.Error())
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return errors.Wrap(ErrInvalidRequest, "HTTP Authorization header missing or invalid")
	}

	client, err := f.Store.GetClient(clientID)
	if err != nil {
		return errors.Wrap(ErrInvalidClient, err.Error())
	}

	// Enforce client authentication for confidential clients
	if !client.IsPublic() {
		if err := f.Hasher.Compare(client.GetHashedSecret(), []byte(clientSecret)); err != nil {
			return errors.Wrap(ErrInvalidClient, err.Error())
		}
	}

	token := r.PostForm.Get("token")
	tokenTypeHint := TokenType(r.PostForm.Get("token_type_hint"))

	var found bool
	for _, loader := range f.RevocationHandlers {
		if err := loader.RevokeToken(ctx, token, tokenTypeHint); err == nil {
			found = true
		} else if errors.Cause(err) == ErrUnknownRequest {
			// do nothing
		} else if err != nil {
			return err
		}
	}

	if !found {
		return errors.WithStack(ErrInvalidRequest)
	}

	return nil
}

// WriteRevocationResponse writes a token revocation response as specified in:
// https://tools.ietf.org/html/rfc7009#section-2.2
//
// The authorization server responds with HTTP status code 200 if the
// token has been revoked successfully or if the client submitted an
// invalid token.
//
// Note: invalid tokens do not cause an error response since the client
// cannot handle such an error in a reasonable way.  Moreover, the
// purpose of the revocation request, invalidating the particular token,
// is already achieved.
func (f *Fosite) WriteRevocationResponse(rw http.ResponseWriter, err error) {
	switch errors.Cause(err) {
	case ErrInvalidRequest:
		fallthrough
	case ErrInvalidClient:
		rw.Header().Set("Content-Type", "application/json;charset=UTF-8")

		rfcerr := ErrorToRFC6749Error(err)
		js, err := json.Marshal(rfcerr)
		if err != nil {
			http.Error(rw, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(rfcerr.StatusCode)
		rw.Write(js)
	default:
		// 200 OK
		rw.WriteHeader(http.StatusOK)
	}
}
