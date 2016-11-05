package warden

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
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

func (w *LocalWarden) IsAllowed(ctx context.Context, a *firewall.AccessRequest) error {
	if err := w.Warden.IsAllowed(&ladon.Request{
		Resource: a.Resource,
		Action:   a.Action,
		Subject:  a.Subject,
		Context:  a.Context,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"subject": a.Subject,
			"request": a,
			"reason":  "The policy decision point denied the request",
		}).WithError(err).Infof("Access denied")
		return err
	}

	return nil
}

func (w *LocalWarden) TokenAllowed(ctx context.Context, token string, a *firewall.TokenAccessRequest, scopes ...string) (*firewall.Context, error) {
	var session = oauth2.NewSession("")
	var auth, err = w.OAuth2.IntrospectToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"subject": session.Subject,
			"request": a,
			"reason":  "Token is expired, malformed or missing",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	return w.sessionAllowed(ctx, a, scopes, auth, session)
}

func (w *LocalWarden) sessionAllowed(ctx context.Context, a *firewall.TokenAccessRequest, scopes []string, oauthRequest fosite.AccessRequester, session *oauth2.Session) (*firewall.Context, error) {
	session = oauthRequest.GetSession().(*oauth2.Session)

	if err := w.Warden.IsAllowed(&ladon.Request{
		Resource: a.Resource,
		Action:   a.Action,
		Subject:  session.Subject,
		Context:  a.Context,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"request":  a,
			"reason":   "The policy decision point denied the request",
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
		ExpiresAt:     session.GetExpiresAt(fosite.AccessToken),
		Extra:         session.Extra,
	}

	logrus.WithFields(logrus.Fields{
		"subject":  c.Subject,
		"audience": oauthRequest.GetClient().GetID(),
	}).Infof("Access granted")

	return c
}
