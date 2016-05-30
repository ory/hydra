package warden

import (
	"net/http"

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

	Issuer string
}

func (w *LocalWarden) actionAllowed(ctx context.Context, a *ladon.Request, scopes []string, oauthRequest fosite.AccessRequester, session *oauth2.Session) (*Context, error) {
	session = oauthRequest.GetSession().(*oauth2.Session)
	if a.Subject != "" && a.Subject != session.Subject {
		return nil, errors.New("Subject mismatch " + a.Subject + " - " + session.Subject)
	}

	if !matchScopes(oauthRequest.GetGrantedScopes(), scopes, session, oauthRequest.GetClient()) {
		return nil, errors.New(herodot.ErrForbidden)
	}

	a.Subject = session.Subject
	if err := w.Warden.IsAllowed(a); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"scopes":   scopes,
		"subject":  a.Subject,
		"audience": oauthRequest.GetClient().GetID(),
		"request":  a,
	}).Infof("Access granted")

	return &Context{
		Subject:       session.Subject,
		GrantedScopes: oauthRequest.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      oauthRequest.GetClient().GetID(),
		IssuedAt:      oauthRequest.GetRequestedAt(),
	}, nil
}

func (w *LocalWarden) ActionAllowed(ctx context.Context, token string, a *ladon.Request, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)
	if err := w.TokenValidator.ValidateToken(ctx, oauthRequest, token); err != nil {
		return nil, err
	}

	return w.actionAllowed(ctx, a, scopes, oauthRequest, session)
}

func (w *LocalWarden) HTTPActionAllowed(ctx context.Context, r *http.Request, a *ladon.Request, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	if err := w.TokenValidator.ValidateRequest(ctx, r, oauthRequest); errors.Is(err, fosite.ErrUnknownRequest) {
		return nil, pkg.ErrUnauthorized
	} else if err != nil {
		return nil, err
	}

	return w.actionAllowed(ctx, a, scopes, oauthRequest, session)
}

func (w *LocalWarden) Authorized(ctx context.Context, token string, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	if err := w.TokenValidator.ValidateToken(ctx, oauthRequest, token); err != nil {
		return nil, err
	}

	session = oauthRequest.GetSession().(*oauth2.Session)
	if !matchScopes(oauthRequest.GetGrantedScopes(), scopes, session, oauthRequest.Client) {
		return nil, errors.New(herodot.ErrForbidden)
	}

	return &Context{
		Subject:       session.Subject,
		GrantedScopes: oauthRequest.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      oauthRequest.GetClient().GetID(),
		IssuedAt:      oauthRequest.GetRequestedAt(),
	}, nil
}

func (w *LocalWarden) HTTPAuthorized(ctx context.Context, r *http.Request, scopes ...string) (*Context, error) {
	var session = new(oauth2.Session)
	var oauthRequest = fosite.NewAccessRequest(session)

	if err := w.TokenValidator.ValidateRequest(ctx, r, oauthRequest); errors.Is(err, fosite.ErrUnknownRequest) {
		return nil, fosite.ErrRequestUnauthorized
	} else if err != nil {
		return nil, err
	}

	session = oauthRequest.GetSession().(*oauth2.Session)
	if !matchScopes(oauthRequest.GetGrantedScopes(), scopes, session, oauthRequest.Client) {
		return nil, errors.New(herodot.ErrForbidden)
	}

	return &Context{
		Subject:       session.Subject,
		GrantedScopes: oauthRequest.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      oauthRequest.GetClient().GetID(),
		IssuedAt:      oauthRequest.GetRequestedAt(),
	}, nil
}

func matchScopes(granted []string, requested []string, session *oauth2.Session, c fosite.Client) bool {
	scopes := &fosite.DefaultScopes{Scopes: granted}
	for _, r := range requested {
		if !scopes.Grant(r) {
			logrus.WithFields(logrus.Fields{
				"reason":           "scope mismatch",
				"granted_scopes":   granted,
				"requested_scopes": requested,
				"audience":         c.GetID(),
				"subject":          session.Subject,
			}).Infof("Authentication failed.")
			return false
		}
	}

	return true
}
