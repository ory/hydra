// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package warden

import (
	"context"
	"net/http"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/firewall"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LocalWarden struct {
	Warden ladon.Warden
	OAuth2 fosite.OAuth2Provider
	Groups group.Manager

	AccessTokenLifespan time.Duration
	Issuer              string
	L                   logrus.FieldLogger
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
		w.L.WithFields(logrus.Fields{
			"subject": a.Subject,
			"request": a,
			"reason":  "The policy decision point denied the request",
		}).WithError(err).Infof("Access denied")
		return err
	}

	w.L.WithFields(logrus.Fields{
		"subject": a.Subject,
		"request": a,
		"reason":  "The policy decision point allowed the request",
	}).Infof("Access allowed")
	return nil
}

func (w *LocalWarden) TokenAllowed(ctx context.Context, token string, a *firewall.TokenAccessRequest, scopes ...string) (*firewall.Context, error) {
	var auth, err = w.OAuth2.IntrospectToken(ctx, token, fosite.AccessToken, oauth2.NewSession(""), scopes...)
	if err != nil {
		w.L.WithFields(logrus.Fields{
			"request": a,
			"scopes":  scopes,
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
		w.L.WithFields(logrus.Fields{
			"scopes":    scopes,
			"subject":   session.GetSubject(),
			"client_id": auth.GetClient().GetID(),
			"request":   a,
			"reason":    "The policy decision point denied the request",
		}).WithError(err).Infof("Access denied")
		return nil, err
	}

	c := w.newContext(auth)
	w.L.WithFields(logrus.Fields{
		"subject":   c.Subject,
		"client_id": auth.GetClient().GetID(),
		"request":   auth,
		"result":    c,
	}).Infof("Access granted")

	return c, nil
}

func (w *LocalWarden) isAllowed(ctx context.Context, a *ladon.Request) error {
	groups, err := w.Groups.FindGroupsByMember(a.Subject)
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
			Subject:  g.ID,
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
		ClientID:      auth.GetClient().GetID(),
		IssuedAt:      auth.GetRequestedAt(),
		ExpiresAt:     exp,
		Extra:         session.Extra,
	}

	return c
}
