package client

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	. "github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

type LocalWarden struct {
	Warden         ladon.Warden
	TokenValidator *core.CoreValidator

	Issuer string
}

func (w *LocalWarden) actionAllowed(ctx context.Context, a *ladon.Request, scopes []string, oauthRequest fosite.AccessRequester, session *oauth2.Session) (*Context, error) {
	if !matchScopes(oauthRequest.GetGrantedScopes(), scopes) {
		return nil, errors.New(herodot.ErrForbidden)
	}

	session = oauthRequest.GetSession().(*oauth2.Session)
	if a.Subject != "" && a.Subject != session.Subject {
		return nil, errors.New("Subject mismatch " + a.Subject + " - " + session.Subject)
	}

	a.Subject = session.Subject
	if err := w.Warden.IsAllowed(a); err != nil {
		return nil, err
	}

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

	if err := w.TokenValidator.ValidateRequest(ctx, r, oauthRequest); err != nil {
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

	if !matchScopes(oauthRequest.GetGrantedScopes(), scopes) {
		return nil, errors.New(herodot.ErrForbidden)
	}

	session = oauthRequest.GetSession().(*oauth2.Session)
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

	if err := w.TokenValidator.ValidateRequest(ctx, r, oauthRequest); err != nil {
		return nil, err
	}

	if !matchScopes(oauthRequest.GetGrantedScopes(), scopes) {
		return nil, errors.New(herodot.ErrForbidden)
	}

	session = oauthRequest.GetSession().(*oauth2.Session)
	return &Context{
		Subject:       session.Subject,
		GrantedScopes: oauthRequest.GetGrantedScopes(),
		Issuer:        w.Issuer,
		Audience:      oauthRequest.GetClient().GetID(),
		IssuedAt:      oauthRequest.GetRequestedAt(),
	}, nil
}

func matchScopes(granted []string, requested []string) bool {
	ga := fosite.Arguments(granted)
	return ga.Has(requested...)
}
