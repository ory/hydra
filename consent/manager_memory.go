/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/pkg"
	"github.com/ory/x/pagination"
)

type MemoryManager struct {
	consentRequests        map[string]ConsentRequest
	handledConsentRequests map[string]HandledConsentRequest
	authRequests           map[string]AuthenticationRequest
	handledAuthRequests    map[string]HandledAuthenticationRequest
	authSessions           map[string]AuthenticationSession
	pairwise               []ForcedObfuscatedAuthenticationSession
	m                      map[string]*sync.RWMutex
	store                  pkg.FositeStorer
}

func NewMemoryManager(store pkg.FositeStorer) *MemoryManager {
	return &MemoryManager{
		consentRequests:        map[string]ConsentRequest{},
		handledConsentRequests: map[string]HandledConsentRequest{},
		authRequests:           map[string]AuthenticationRequest{},
		handledAuthRequests:    map[string]HandledAuthenticationRequest{},
		authSessions:           map[string]AuthenticationSession{},
		pairwise:               []ForcedObfuscatedAuthenticationSession{},
		store:                  store,
		m: map[string]*sync.RWMutex{
			"consentRequests":        new(sync.RWMutex),
			"handledConsentRequests": new(sync.RWMutex),
			"authRequests":           new(sync.RWMutex),
			"handledAuthRequests":    new(sync.RWMutex),
			"authSessions":           new(sync.RWMutex),
		},
	}
}

func (m *MemoryManager) CreateForcedObfuscatedAuthenticationSession(ctx context.Context, s *ForcedObfuscatedAuthenticationSession) error {
	for k, v := range m.pairwise {
		if v.Subject == s.Subject && v.ClientID == s.ClientID {
			m.pairwise[k] = *s
			return nil
		}
	}

	m.pairwise = append(m.pairwise, *s)
	return nil
}

func (m *MemoryManager) GetForcedObfuscatedAuthenticationSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedAuthenticationSession, error) {
	for _, v := range m.pairwise {
		if v.SubjectObfuscated == obfuscated && v.ClientID == client {
			return &v, nil
		}
	}

	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) RevokeUserConsentSession(ctx context.Context, user string) error {
	return m.RevokeUserClientConsentSession(ctx, user, "")
}

func (m *MemoryManager) RevokeUserClientConsentSession(ctx context.Context, user, client string) error {
	m.m["handledConsentRequests"].Lock()
	defer m.m["handledConsentRequests"].Unlock()

	var found bool
	for k, c := range m.handledConsentRequests {
		cr, err := m.GetConsentRequest(ctx, c.Challenge)
		if err != nil {
			return err
		}

		if cr.Subject == user &&
			(client == "" ||
				(client != "" && cr.Client.GetID() == client)) {
			delete(m.handledConsentRequests, k)

			m.m["consentRequests"].Lock()
			delete(m.consentRequests, k)
			m.m["consentRequests"].Unlock()

			if err := m.store.RevokeAccessToken(nil, c.Challenge); errors.Cause(err) == fosite.ErrNotFound {
				// do nothing
			} else if err != nil {
				return err
			}
			if err := m.store.RevokeRefreshToken(nil, c.Challenge); errors.Cause(err) == fosite.ErrNotFound {
				// do nothing
			} else if err != nil {
				return err
			}
			found = true
		}
	}

	if !found {
		return errors.WithStack(pkg.ErrNotFound)
	}
	return nil
}

func (m *MemoryManager) RevokeUserAuthenticationSession(ctx context.Context, user string) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()

	var found bool
	for k, c := range m.authSessions {
		if c.Subject == user {
			delete(m.authSessions, k)
			found = true
		}
	}

	if !found {
		return errors.WithStack(pkg.ErrNotFound)
	}
	return nil
}

func (m *MemoryManager) CreateConsentRequest(ctx context.Context, c *ConsentRequest) error {
	m.m["consentRequests"].Lock()
	defer m.m["consentRequests"].Unlock()
	if _, ok := m.consentRequests[c.Challenge]; ok {
		return errors.New("Key already exists")
	}
	m.consentRequests[c.Challenge] = *c
	return nil
}

func (m *MemoryManager) GetConsentRequest(ctx context.Context, challenge string) (*ConsentRequest, error) {
	m.m["consentRequests"].RLock()
	defer m.m["consentRequests"].RUnlock()

	c, ok := m.consentRequests[challenge]
	if !ok {
		return nil, errors.WithStack(pkg.ErrNotFound)
	}

	for _, h := range m.handledConsentRequests {
		if h.Challenge == c.Challenge {
			c.WasHandled = h.WasUsed
		}
	}
	c.Client.ClientID = c.Client.GetID()
	return &c, nil
}

func (m *MemoryManager) HandleConsentRequest(ctx context.Context, challenge string, r *HandledConsentRequest) (*ConsentRequest, error) {
	m.m["handledConsentRequests"].Lock()
	m.handledConsentRequests[r.Challenge] = *r
	m.m["handledConsentRequests"].Unlock()
	return m.GetConsentRequest(ctx, challenge)
}

func (m *MemoryManager) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*HandledConsentRequest, error) {
	for _, c := range m.consentRequests {
		if c.Verifier == verifier {
			for _, h := range m.handledConsentRequests {
				if h.Challenge == c.Challenge {
					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleConsentRequest(ctx, h.Challenge, &h); err != nil {
						return nil, err
					}

					c.Client.ClientID = c.Client.GetID()
					h.ConsentRequest = &c
					return &h, nil
				}
			}
		}
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) FindPreviouslyGrantedConsentRequests(ctx context.Context, client, subject string) ([]HandledConsentRequest, error) {
	var rs []HandledConsentRequest
	filteredByUser, err := m.FindPreviouslyGrantedConsentRequestsByUser(ctx, subject, -1, -1)
	if errors.Cause(err) == pkg.ErrNotFound {
		return nil, errors.WithStack(ErrNoPreviousConsentFound)
	} else if err != nil {
		return nil, err
	}

	for _, c := range filteredByUser {
		if client == c.ConsentRequest.Client.GetID() {
			rs = append(rs, c)
		}
	}

	if len(rs) == 0 {
		return nil, errors.WithStack(ErrNoPreviousConsentFound)
	}

	return rs, nil
}

func (m *MemoryManager) FindPreviouslyGrantedConsentRequestsByUser(ctx context.Context, subject string, limit, offset int) ([]HandledConsentRequest, error) {
	var rs []HandledConsentRequest
	for _, c := range m.handledConsentRequests {
		cr, err := m.GetConsentRequest(ctx, c.Challenge)
		if err != nil {
			return nil, err
		}

		if subject != cr.Subject {
			continue
		}

		if c.Error != nil {
			continue
		}

		if !c.Remember {
			continue
		}

		if cr.Skip {
			continue
		}

		if c.RememberFor > 0 && c.RequestedAt.Add(time.Duration(c.RememberFor)*time.Second).Before(time.Now().UTC()) {
			continue
		}

		cr.Client.ClientID = cr.Client.GetID()
		c.ConsentRequest = cr
		rs = append(rs, c)
	}

	if len(rs) == 0 {
		return nil, errors.WithStack(ErrNoPreviousConsentFound)
	}

	if limit < 0 && offset < 0 {
		return rs, nil
	}

	start, end := pagination.Index(limit, offset, len(rs))
	return rs[start:end], nil
}

func (m *MemoryManager) GetAuthenticationSession(ctx context.Context, id string) (*AuthenticationSession, error) {
	m.m["authSessions"].RLock()
	defer m.m["authSessions"].RUnlock()
	if c, ok := m.authSessions[id]; ok {
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) CreateAuthenticationSession(ctx context.Context, a *AuthenticationSession) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	if _, ok := m.authSessions[a.ID]; ok {
		return errors.New("Key already exists")
	}
	m.authSessions[a.ID] = *a
	return nil
}

func (m *MemoryManager) DeleteAuthenticationSession(ctx context.Context, id string) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	delete(m.authSessions, id)
	return nil
}

func (m *MemoryManager) CreateAuthenticationRequest(ctx context.Context, a *AuthenticationRequest) error {
	m.m["authRequests"].Lock()
	defer m.m["authRequests"].Unlock()
	if _, ok := m.authRequests[a.Challenge]; ok {
		return errors.New("Key already exists")
	}
	m.authRequests[a.Challenge] = *a
	return nil
}

func (m *MemoryManager) GetAuthenticationRequest(ctx context.Context, challenge string) (*AuthenticationRequest, error) {
	m.m["authRequests"].RLock()
	defer m.m["authRequests"].RUnlock()

	c, ok := m.authRequests[challenge]
	if !ok {
		return nil, errors.WithStack(pkg.ErrNotFound)
	}

	for _, h := range m.handledAuthRequests {
		if h.Challenge == c.Challenge {
			c.WasHandled = h.WasUsed
		}
	}
	c.Client.ClientID = c.Client.GetID()
	return &c, nil
}

func (m *MemoryManager) HandleAuthenticationRequest(ctx context.Context, challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error) {
	m.m["handledAuthRequests"].Lock()
	m.handledAuthRequests[r.Challenge] = *r
	m.m["handledAuthRequests"].Unlock()
	return m.GetAuthenticationRequest(ctx, challenge)
}

func (m *MemoryManager) VerifyAndInvalidateAuthenticationRequest(ctx context.Context, verifier string) (*HandledAuthenticationRequest, error) {
	for _, c := range m.authRequests {
		if c.Verifier == verifier {
			for _, h := range m.handledAuthRequests {
				if h.Challenge == c.Challenge {
					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Authentication verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleAuthenticationRequest(ctx, h.Challenge, &h); err != nil {
						return nil, err
					}

					c.Client.ClientID = c.Client.GetID()
					h.AuthenticationRequest = &c
					return &h, nil
				}
			}
		}
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}
