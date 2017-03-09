package fosite

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type TokenIntrospector interface {
	IntrospectToken(ctx context.Context, token string, tokenType TokenType, accessRequest AccessRequester, scopes []string) error
}

func AccessTokenFromRequest(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		return ""
	}

	return split[1]
}

func (f *Fosite) IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scopes ...string) (AccessRequester, error) {
	var found bool = false

	ar := NewAccessRequest(session)
	for _, validator := range f.TokenIntrospectionHandlers {
		if err := errors.Cause(validator.IntrospectToken(ctx, token, tokenType, ar, scopes)); err == ErrUnknownRequest {
			// Nothing to do
		} else if err != nil {
			return nil, errors.Wrap(err, "A validator returned an error")
		} else {
			found = true
		}
	}

	if !found {
		return nil, errors.Wrap(ErrRequestUnauthorized, "No validator felt responsible for validating the token")
	}

	return ar, nil
}
