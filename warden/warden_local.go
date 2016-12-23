package warden

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/warden/group"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

type LocalWarden struct {
	Warden ladon.Warden
	OAuth2 fosite.OAuth2Provider
	Groups group.Manager

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
	var auth, err = w.OAuth2.IntrospectToken(ctx, token, fosite.AccessToken, oauth2.NewSession(""), scopes...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": a,
			"reason":  "Token is expired, malformed or missing",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	session := auth.GetSession()
	c, err := w.sessionAllowed(ctx, a, scopes, auth, session.GetSubject())
	if err != nil {
		orig := err
		// If the subject is not allowed, check if a group he belongs to, is.
		groups, err := w.Groups.FindGroupNames(session.GetSubject())
		if err != nil {
			return nil, err
		} else if len(groups) == 0 {
			return nil, orig
		}

		logrus.WithFields(logrus.Fields{
			"subject": session.GetSubject(),
			"request": a,
			"groups":  groups,
			"reason":  "Subject is not allowed to perform action on resource, trying groups",
		}).WithError(err).Infof("Access denied")
		for _, group := range groups {
			// If one of the groups allows access, allow access.
			if c, err := w.sessionAllowed(ctx, a, scopes, auth, group); err == nil {
				logrus.WithFields(logrus.Fields{
					"subject":  c.Subject,
					"group":    group,
					"audience": auth.GetClient().GetID(),
				}).Infof("Access granted, because subject is member of authorized group")
				return c, nil
			}

			// We don't really care about errors here
		}
		return c, orig
	}

	return c, err
}

func (w *LocalWarden) sessionAllowed(ctx context.Context, a *firewall.TokenAccessRequest, scopes []string, oauthRequest fosite.AccessRequester, subject string) (*firewall.Context, error) {
	if err := w.Warden.IsAllowed(&ladon.Request{
		Resource: a.Resource,
		Action:   a.Action,
		Subject:  subject,
		Context:  a.Context,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  subject,
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

	exp := oauthRequest.GetSession().GetExpiresAt(fosite.AccessToken)
	if exp.IsZero() {
		exp = oauthRequest.GetRequestedAt().Add(w.AccessTokenLifespan)
	}
	c := &firewall.Context{
		Subject:       session.Subject,
		GrantedScopes: oauthRequest.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      oauthRequest.GetClient().GetID(),
		IssuedAt:      oauthRequest.GetRequestedAt(),
		ExpiresAt:     exp,
		Extra:         session.Extra,
	}

	logrus.WithFields(logrus.Fields{
		"subject":  c.Subject,
		"audience": oauthRequest.GetClient().GetID(),
	}).Infof("Access granted")

	return c
}
