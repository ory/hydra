package oauth2

import (
	"github.com/Sirupsen/logrus"
	"github.com/ory-am/fosite"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

type LocalIntrospector struct {
	OAuth2 fosite.OAuth2Provider

	AccessTokenLifespan time.Duration
	Issuer              string
}

func (w *LocalIntrospector) TokenFromRequest(r *http.Request) string {
	return fosite.AccessTokenFromRequest(r)
}

func (w *LocalIntrospector) IntrospectToken(ctx context.Context, token string, scopes ...string) (*Introspection, error) {
	var session = new(Session)
	var auth, err = w.OAuth2.ValidateToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logrus.WithError(err).Infof("Token introspection failed")
		return &Introspection{
			Active: false,
		}, err
	}

	session = auth.GetSession().(*Session)
	return &Introspection{
		Active:    true,
		Subject:   session.Subject,
		Audience:  auth.GetClient().GetID(),
		Scope:     strings.Join(auth.GetGrantedScopes(), " "),
		Issuer:    w.Issuer,
		IssuedAt:  auth.GetRequestedAt().Unix(),
		NotBefore: auth.GetRequestedAt().Unix(),
		ExpiresAt: session.AccessTokenExpiresAt(auth.GetRequestedAt().Add(w.AccessTokenLifespan)).Unix(),
		Extra:     session.Extra,
	}, nil
}
