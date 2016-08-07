package warden

import (
	"net/http"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

type LocalWarden struct {
	Warden              ladon.Warden
	OAuth2              fosite.OAuth2Provider

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
			"reason":  "request denied by policies",
		}).Infof("Access denied")
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
			"reason":  "token could not be validated",
		}).Infof("Access denied")
		return nil, err
	}

	return w.allowed(ctx, a, scopes, auth, session)
}

func (w *LocalWarden) InspectToken(ctx context.Context, token string, scopes ...string) (*firewall.Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	var auth, err = w.OAuth2.ValidateToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"reason":   "token validation failed",
		}).Infof("Access denied")
		return nil, err
	}

	return w.newContext(auth), nil
}

func (w *LocalWarden) allowed(ctx context.Context, a *ladon.Request, scopes []string, oauthRequest fosite.AccessRequester, session *oauth2.Session) (*firewall.Context, error) {
	session = oauthRequest.GetSession().(*oauth2.Session)
	if a.Subject != "" && a.Subject != session.Subject {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  a.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"request":  a,
			"reason":   "subject mismatch",
		}).Infof("Access denied")
		return nil, errors.Errorf("Expected subject to be %s but got %s", session.Subject, a.Subject)
	}

	a.Subject = session.Subject
	if err := w.Warden.IsAllowed(a); err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  a.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"request":  a,
			"reason":   "policy effect is deny",
		}).Infof("Access denied")
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
