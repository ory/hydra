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
	"github.com/pkg/errors"
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
	if err := w.isAllowed(ctx, &ladon.Request{
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

	logrus.WithFields(logrus.Fields{
		"subject": a.Subject,
		"request": a,
		"reason":  "The policy decision point allowed the request",
	}).Infof("Access allowed")
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
	if err := w.isAllowed(ctx, &ladon.Request{
		Resource: a.Resource,
		Action:   a.Action,
		Subject:  session.GetSubject(),
		Context:  a.Context,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.GetSubject(),
			"audience": auth.GetClient().GetID(),
			"request":  a,
			"reason":   "The policy decision point denied the request",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	c := w.newContext(auth)
	logrus.WithFields(logrus.Fields{
		"subject":  c.Subject,
		"audience": auth.GetClient().GetID(),
		"request":  auth,
		"result":   c,
	}).Infof("Access granted")

	return c, nil
}

func (w *LocalWarden) isAllowed(ctx context.Context, a *ladon.Request) error {
	groups, err := w.Groups.FindGroupNames(a.Subject)
	if err != nil {
		return err
	}

	errs := make([]error, len(groups)+1)
	errs[0] = w.Warden.IsAllowed(&ladon.Request{
		Resource: a.Resource,
		Action:   a.Action,
		Subject:  a.Subject,
		Context:  a.Context,
	})

	for k, g := range groups {
		errs[k+1] = w.Warden.IsAllowed(&ladon.Request{
			Resource: a.Resource,
			Action:   a.Action,
			Subject:  g,
			Context:  a.Context,
		})
	}

	for _, err := range errs {
		if errors.Cause(err) == ladon.ErrRequestForcefullyDenied {
			return errors.Wrap(fosite.ErrRequestForbidden, err.Error())
		}
	}

	for _, err := range errs {
		if err == nil {
			return nil
		}
	}

	return errors.Wrap(fosite.ErrRequestForbidden, ladon.ErrRequestDenied.Error())
}

func (w *LocalWarden) newContext(auth fosite.AccessRequester) *firewall.Context {
	session := auth.GetSession().(*oauth2.Session)

	exp := auth.GetSession().GetExpiresAt(fosite.AccessToken)
	if exp.IsZero() {
		exp = auth.GetRequestedAt().Add(w.AccessTokenLifespan)
	}

	c := &firewall.Context{
		Subject:       session.Subject,
		GrantedScopes: auth.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      auth.GetClient().GetID(),
		IssuedAt:      auth.GetRequestedAt(),
		ExpiresAt:     exp,
		Extra:         session.Extra,
	}

	return c
}
