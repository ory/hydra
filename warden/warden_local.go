package warden

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

type LocalWarden struct {
	Warden ladon.Warden
	OAuth2 fosite.OAuth2Provider

	AccessTokenLifespan time.Duration
	Issuer              string
}

func (w *LocalWarden) TokenFromRequest(r *http.Request) string {
	return fosite.AccessTokenFromRequest(r)
}

func (w *LocalWarden) IsAllowed(ctx context.Context, a *ladon.Request) error {
	if err := w.Warden.IsAllowed(a); err != nil {
		logrus.WithFields(logrus.Fields{
			"subject": a.Subject,
			"request": a,
			"reason":  "The policy decision point denied the request",
		}).WithError(err).Infof("Access denied")
		return err
	}

	return nil
}

func (w *LocalWarden) TokenAllowed(ctx context.Context, token string, a *ladon.Request, scopes ...string) (*firewall.Context, error) {
	var session = new(oauth2.Session)
	var auth, err = w.OAuth2.ValidateToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"subject": a.Subject,
			"request": a,
			"reason":  "Token is expired, malformed or missing",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	return w.sessionAllowed(ctx, a, scopes, auth, session)
}

func (w *LocalWarden) TokenValid(ctx context.Context, token string, scopes ...string) (*firewall.Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	var auth, err = w.OAuth2.ValidateToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"reason":  "Token is expired, malformed or missing",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	return w.newContext(auth), nil
}

func (w *LocalWarden) sessionAllowed(ctx context.Context, a *ladon.Request, scopes []string, oauthRequest fosite.AccessRequester, session *oauth2.Session) (*firewall.Context, error) {
	session = oauthRequest.GetSession().(*oauth2.Session)
	if a.Subject != "" && a.Subject != session.Subject {
		err := errors.Errorf("Expected subject to be %s but got %s", session.Subject, a.Subject)
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  a.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"request":  a,
			"reason":   "Request subject and token subject do not match",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	a.Subject = session.Subject
	if err := w.Warden.IsAllowed(a); err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  a.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"request":  a,
			"reason":  "The policy decision point denied the request",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	return w.newContext(oauthRequest), nil
}

func (w *LocalWarden) newContext(oauthRequest fosite.AccessRequester) *firewall.Context {
	session := oauthRequest.GetSession().(*oauth2.Session)
	c := &firewall.Context{
		Subject:       session.Subject,
		GrantedScopes: oauthRequest.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      oauthRequest.GetClient().GetID(),
		IssuedAt:      oauthRequest.GetRequestedAt(),
		ExpiresAt:     session.AccessTokenExpiresAt(oauthRequest.GetRequestedAt().Add(w.AccessTokenLifespan)),
		Extra:         session.Extra,
	}

	logrus.WithFields(logrus.Fields{
		"subject":  c.Subject,
		"audience": oauthRequest.GetClient().GetID(),
	}).Infof("Access granted")

	return c
}
