package warden

import (
	"net/http"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	. "github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

type LocalWarden struct {
	Warden         ladon.Warden
	TokenValidator *core.CoreValidator

	AccessTokenLifespan time.Duration
	Issuer              string
}

func (w *LocalWarden) IsAllowed(ctx context.Context, a *ladon.Request) error {
	if err := w.Warden.IsAllowed(a); err != nil {
		logrus.WithFields(logrus.Fields{
			"subject": a.Subject,
			"request": a,
			"reason":  "policy effect is deny",
		}).Infof("Access denied")
		return err
	}

	return nil
}

func (w *LocalWarden) TokenAllowed(ctx context.Context, token string, a *ladon.Request, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)
	if err := w.TokenValidator.ValidateToken(ctx, oauthRequest, token); err != nil {
		return nil, err
	}

	return w.allowed(ctx, a, scopes, oauthRequest, session)
}

func (w *LocalWarden) HTTPRequestAllowed(ctx context.Context, r *http.Request, a *ladon.Request, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	if err := w.validateRequest(ctx, r, oauthRequest, scopes); err != nil {
		return nil, err
	}

	return w.allowed(ctx, a, scopes, oauthRequest, session)
}

func (w *LocalWarden) InspectToken(ctx context.Context, token string, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	if err := w.validateToken(ctx, token, scopes, oauthRequest); err != nil {
		return nil, err
	}

	if err := matchScopes(oauthRequest, scopes); err != nil {
		return nil, err
	}

	return w.newContext(oauthRequest), nil
}

func (w *LocalWarden) InspectTokenFromHTTP(ctx context.Context, r *http.Request, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	if err := w.validateRequest(ctx, r, oauthRequest, scopes); err != nil {
		return nil, err
	}

	if err := matchScopes(oauthRequest, scopes); err != nil {
		return nil, err
	}

	return w.newContext(oauthRequest), nil
}

func (w *LocalWarden) allowed(ctx context.Context, a *ladon.Request, scopes []string, oauthRequest fosite.AccessRequester, session *oauth2.Session) (*Context, error) {
	session = oauthRequest.GetSession().(*oauth2.Session)
	if a.Subject != "" && a.Subject != session.Subject {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  a.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"request":  a,
			"reason":   "subject mismatch",
		}).Infof("Access denied")
		return nil, errors.New("Subject mismatch " + a.Subject + " - " + session.Subject)
	}

	if err := matchScopes(oauthRequest, scopes); err != nil {
		return nil, err
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

func (w *LocalWarden) validateToken(ctx context.Context, token string, scopes []string, oauthRequest *fosite.AccessRequest) error {
	session := oauthRequest.GetSession().(*oauth2.Session)
	if err := w.TokenValidator.ValidateToken(ctx, oauthRequest, token); err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"reason":   "token validation failed",
		}).Infof("Access denied")
		return err
	}
	return nil
}

func (w *LocalWarden) validateRequest(ctx context.Context, r *http.Request, oauthRequest *fosite.AccessRequest, scopes []string) error {
	session := oauthRequest.GetSession().(*oauth2.Session)
	if err := w.TokenValidator.ValidateRequest(ctx, r, oauthRequest); errors.Is(err, fosite.ErrUnknownRequest) {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"reason":   "unknown request",
		}).Infof("Access denied")
		return errors.New(pkg.ErrUnauthorized)
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"scopes":   scopes,
			"subject":  session.Subject,
			"audience": oauthRequest.GetClient().GetID(),
			"reason":   "token validation failed",
		}).Infof("Access denied")
		return err
	}
	return nil
}

func (w *LocalWarden) newContext(oauthRequest fosite.AccessRequester) *Context {
	session := oauthRequest.GetSession().(*oauth2.Session)
	c := &Context{
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

func matchScopes(oauthRequest fosite.AccessRequester, requested []string) error {
	session := oauthRequest.GetSession().(*oauth2.Session)
	scopes := &fosite.DefaultScopes{Scopes: oauthRequest.GetGrantedScopes()}
	for _, r := range requested {
		if !scopes.Grant(r) {
			logrus.WithFields(logrus.Fields{
				"reason":           "Scopes are not matching",
				"granted_scopes":   oauthRequest.GetGrantedScopes,
				"requested_scopes": requested,
				"audience":         oauthRequest.GetClient().GetID(),
				"subject":          session.Subject,
			}).Infof("Access denied.")
			return errors.New(herodot.ErrForbidden)
		}
	}

	return nil
}
