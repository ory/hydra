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

	"github.com/ory/hydra/client"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/x/pagination"

	"github.com/ory/hydra/x"
)

type MemoryManager struct {
	consentRequests        map[string]ConsentRequest
	logoutRequests         map[string]LogoutRequest
	handledConsentRequests map[string]HandledConsentRequest
	authRequests           map[string]LoginRequest
	handledAuthRequests    map[string]HandledLoginRequest
	authSessions           map[string]LoginSession
	pairwise               []ForcedObfuscatedLoginSession
	m                      map[string]*sync.RWMutex
	r                      InternalRegistry
}

func NewMemoryManager(r InternalRegistry) *MemoryManager {
	return &MemoryManager{
		consentRequests:        map[string]ConsentRequest{},
		logoutRequests:         map[string]LogoutRequest{},
		handledConsentRequests: map[string]HandledConsentRequest{},
		authRequests:           map[string]LoginRequest{},
		handledAuthRequests:    map[string]HandledLoginRequest{},
		authSessions:           map[string]LoginSession{},
		pairwise:               []ForcedObfuscatedLoginSession{},
		r:                      r,
		m: map[string]*sync.RWMutex{
			"logoutRequests":         new(sync.RWMutex),
			"consentRequests":        new(sync.RWMutex),
			"handledConsentRequests": new(sync.RWMutex),
			"authRequests":           new(sync.RWMutex),
			"handledAuthRequests":    new(sync.RWMutex),
			"authSessions":           new(sync.RWMutex),
		},
	}
}

func (m *MemoryManager) CreateForcedObfuscatedLoginSession(ctx context.Context, s *ForcedObfuscatedLoginSession) error {
	for k, v := range m.pairwise {
		if v.Subject == s.Subject && v.ClientID == s.ClientID {
			m.pairwise[k] = *s
			return nil
		}
	}

	m.pairwise = append(m.pairwise, *s)
	return nil
}

func (m *MemoryManager) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedLoginSession, error) {
	for _, v := range m.pairwise {
		if v.SubjectObfuscated == obfuscated && v.ClientID == client {
			return &v, nil
		}
	}

	return nil, errors.WithStack(x.ErrNotFound)
}

func (m *MemoryManager) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	return m.RevokeSubjectClientConsentSession(ctx, user, "")
}

func (m *MemoryManager) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	var found bool

	m.m["handledConsentRequests"].RLock()
	for k, c := range m.handledConsentRequests {
		m.m["handledConsentRequests"].RUnlock()
		cr, err := m.GetConsentRequest(ctx, c.Challenge)
		if err != nil {
			return err
		}
		m.m["handledConsentRequests"].RLock()

		if cr.Subject == user &&
			(client == "" ||
				(client != "" && cr.Client.GetID() == client)) {
			delete(m.handledConsentRequests, k)

			m.m["consentRequests"].Lock()
			delete(m.consentRequests, k)
			m.m["consentRequests"].Unlock()

			if err := m.r.OAuth2Storage().RevokeAccessToken(ctx, c.Challenge); errors.Cause(err) == fosite.ErrNotFound {
				// do nothing
			} else if err != nil {
				return err
			}
			if err := m.r.OAuth2Storage().RevokeRefreshToken(ctx, c.Challenge); errors.Cause(err) == fosite.ErrNotFound {
				// do nothing
			} else if err != nil {
				return err
			}
			found = true
		}
	}
	m.m["handledConsentRequests"].RUnlock()

	if !found {
		return errors.WithStack(x.ErrNotFound)
	}
	return nil
}

func (m *MemoryManager) RevokeSubjectLoginSession(ctx context.Context, user string) error {
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
		return errors.WithStack(x.ErrNotFound)
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
		return nil, errors.WithStack(x.ErrNotFound)
	}

	m.m["handledConsentRequests"].RLock()
	for _, h := range m.handledConsentRequests {
		if h.Challenge == c.Challenge {
			c.WasHandled = h.WasUsed
		}
	}
	m.m["handledConsentRequests"].RUnlock()

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
	m.m["consentRequests"].RLock()
	for _, c := range m.consentRequests {
		if c.Verifier == verifier {
			m.m["handledConsentRequests"].RLock()
			for _, h := range m.handledConsentRequests {
				if h.Challenge == c.Challenge {
					m.m["consentRequests"].RUnlock()
					m.m["handledConsentRequests"].RUnlock()
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
			m.m["handledConsentRequests"].RUnlock()
		}
	}
	m.m["consentRequests"].RUnlock()
	return nil, errors.WithStack(x.ErrNotFound)
}

func (m *MemoryManager) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, subject string) ([]HandledConsentRequest, error) {
	var rs []HandledConsentRequest

	m.m["handledConsentRequests"].RLock()
	for _, c := range m.handledConsentRequests {
		m.m["handledConsentRequests"].RUnlock()
		cr, err := m.GetConsentRequest(ctx, c.Challenge)
		if errors.Cause(err) == x.ErrNotFound {
			return nil, errors.WithStack(ErrNoPreviousConsentFound)
		} else if err != nil {
			return nil, err
		}
		m.m["handledConsentRequests"].RLock()

		if subject != cr.Subject {
			continue
		}

		if client != cr.Client.GetID() {
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
	m.m["handledConsentRequests"].RUnlock()

	if len(rs) == 0 {
		return nil, errors.WithStack(ErrNoPreviousConsentFound)
	}

	return rs, nil
}

func (m *MemoryManager) FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit, offset int) ([]HandledConsentRequest, error) {
	var rs []HandledConsentRequest

	m.m["handledConsentRequests"].RLock()
	for _, c := range m.handledConsentRequests {
		m.m["handledConsentRequests"].RUnlock()
		cr, err := m.GetConsentRequest(ctx, c.Challenge)
		if err != nil {
			return nil, err
		}
		m.m["handledConsentRequests"].RLock()

		if subject != cr.Subject {
			continue
		}

		if c.Error != nil {
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
	m.m["handledConsentRequests"].RUnlock()

	if len(rs) == 0 {
		return nil, errors.WithStack(ErrNoPreviousConsentFound)
	}

	if limit < 0 && offset < 0 {
		return rs, nil
	}

	start, end := pagination.Index(limit, offset, len(rs))
	return rs[start:end], nil
}

func (m *MemoryManager) CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error) {
	var rs []HandledConsentRequest
	for _, c := range m.handledConsentRequests {
		cr, err := m.GetConsentRequest(ctx, c.Challenge)
		if err != nil {
			return 0, err
		}

		if subject != cr.Subject {
			continue
		}

		if c.Error != nil {
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

	return len(rs), nil
}

func (m *MemoryManager) ConfirmLoginSession(ctx context.Context, id string, subject string, remember bool) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	if c, ok := m.authSessions[id]; ok {
		c.Remember = remember
		c.Subject = subject
		c.AuthenticatedAt = time.Now().UTC()
		m.authSessions[id] = c
		return nil
	}
	return errors.WithStack(x.ErrNotFound)
}

func (m *MemoryManager) GetRememberedLoginSession(ctx context.Context, id string) (*LoginSession, error) {
	m.m["authSessions"].RLock()
	defer m.m["authSessions"].RUnlock()
	if c, ok := m.authSessions[id]; ok {
		if c.Remember {
			return &c, nil
		}
		return nil, errors.WithStack(x.ErrNotFound)
	}
	return nil, errors.WithStack(x.ErrNotFound)
}

func (m *MemoryManager) CreateLoginSession(ctx context.Context, a *LoginSession) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	if _, ok := m.authSessions[a.ID]; ok {
		return errors.New("Key already exists")
	}
	m.authSessions[a.ID] = *a
	return nil
}

func (m *MemoryManager) DeleteLoginSession(ctx context.Context, id string) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	delete(m.authSessions, id)
	return nil
}

func (m *MemoryManager) CreateLoginRequest(ctx context.Context, a *LoginRequest) error {
	m.m["authRequests"].Lock()
	defer m.m["authRequests"].Unlock()
	if _, ok := m.authRequests[a.Challenge]; ok {
		return errors.New("Key already exists")
	}
	m.authRequests[a.Challenge] = *a
	return nil
}

func (m *MemoryManager) GetLoginRequest(ctx context.Context, challenge string) (*LoginRequest, error) {
	m.m["authRequests"].RLock()
	defer m.m["authRequests"].RUnlock()

	c, ok := m.authRequests[challenge]
	if !ok {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	for _, h := range m.handledAuthRequests {
		if h.Challenge == c.Challenge {
			c.WasHandled = h.WasUsed
		}
	}
	c.Client.ClientID = c.Client.GetID()
	return &c, nil
}

func (m *MemoryManager) HandleLoginRequest(ctx context.Context, challenge string, r *HandledLoginRequest) (*LoginRequest, error) {
	m.m["handledAuthRequests"].Lock()
	m.handledAuthRequests[r.Challenge] = *r
	m.m["handledAuthRequests"].Unlock()
	return m.GetLoginRequest(ctx, challenge)
}

func (m *MemoryManager) VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*HandledLoginRequest, error) {
	m.m["authRequests"].RLock()
	for _, c := range m.authRequests {
		if c.Verifier == verifier {
			m.m["handledAuthRequests"].RLock()
			for _, h := range m.handledAuthRequests {
				if h.Challenge == c.Challenge {
					m.m["handledAuthRequests"].RUnlock()
					m.m["authRequests"].RUnlock()

					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Authentication verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleLoginRequest(ctx, h.Challenge, &h); err != nil {
						return nil, err
					}

					c.Client.ClientID = c.Client.GetID()
					h.LoginRequest = &c
					return &h, nil
				}
			}
			m.m["handledAuthRequests"].RUnlock()
		}
	}

	m.m["authRequests"].RUnlock()
	return nil, errors.WithStack(x.ErrNotFound)
}

func (m *MemoryManager) ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	m.m["consentRequests"].RLock()
	defer m.m["consentRequests"].RUnlock()

	preventDupes := make(map[string]bool)
	var rs []client.Client
	for _, cr := range m.consentRequests {
		if cr.Subject == subject &&
			len(cr.Client.FrontChannelLogoutURI) > 0 &&
			cr.LoginSessionID == sid &&
			!preventDupes[cr.Client.GetID()] {

			rs = append(rs, *cr.Client)
			preventDupes[cr.Client.GetID()] = true
		}
	}

	return rs, nil
}

func (m *MemoryManager) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	m.m["consentRequests"].RLock()
	defer m.m["consentRequests"].RUnlock()

	clientsMap := make(map[string]bool)
	var rs []client.Client
	for _, cr := range m.consentRequests {
		if cr.Subject == subject &&
			cr.LoginSessionID == sid &&
			len(cr.Client.BackChannelLogoutURI) > 0 &&
			!(clientsMap[cr.Client.GetID()]) {
			rs = append(rs, *cr.Client)
			clientsMap[cr.Client.GetID()] = true
		}
	}

	return rs, nil
}

func (m *MemoryManager) CreateLogoutRequest(ctx context.Context, r *LogoutRequest) error {
	m.m["logoutRequests"].Lock()
	m.logoutRequests[r.Challenge] = *r
	m.m["logoutRequests"].Unlock()
	return nil
}

func (m *MemoryManager) GetLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	m.m["logoutRequests"].RLock()
	defer m.m["logoutRequests"].RUnlock()
	if c, ok := m.logoutRequests[challenge]; ok {
		return &c, nil
	}
	return nil, errors.WithStack(x.ErrNotFound)
}

func (m *MemoryManager) AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error) {
	m.m["logoutRequests"].Lock()
	lr, ok := m.logoutRequests[challenge]
	if !ok {
		m.m["logoutRequests"].Unlock()
		return nil, errors.WithStack(x.ErrNotFound)
	}

	lr.Accepted = true
	m.logoutRequests[challenge] = lr

	m.m["logoutRequests"].Unlock()
	return m.GetLogoutRequest(ctx, challenge)
}

func (m *MemoryManager) RejectLogoutRequest(ctx context.Context, challenge string) error {
	m.m["logoutRequests"].Lock()
	defer m.m["logoutRequests"].Unlock()

	if _, ok := m.logoutRequests[challenge]; !ok {
		return errors.WithStack(x.ErrNotFound)
	}

	delete(m.logoutRequests, challenge)
	return nil
}

func (m *MemoryManager) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*LogoutRequest, error) {
	m.m["logoutRequests"].RLock()
	for _, c := range m.logoutRequests {
		if c.Verifier == verifier {
			m.m["logoutRequests"].RUnlock()

			if c.WasUsed {
				return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Logout verifier has been used already"))
			}

			if !c.Accepted {
				return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Logout verifier has not been accepted yet"))
			}

			c.WasUsed = true
			if err := m.CreateLogoutRequest(ctx, &c); err != nil {
				return nil, err
			}

			return &c, nil
		}
	}

	m.m["logoutRequests"].RUnlock()
	return nil, errors.WithStack(x.ErrNotFound)
}
