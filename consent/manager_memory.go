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
	"sync"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type MemoryManager struct {
	consentRequests        map[string]ConsentRequest
	handledConsentRequests map[string]HandledConsentRequest
	authRequests           map[string]AuthenticationRequest
	handledAuthRequests    map[string]HandledAuthenticationRequest
	authSessions           map[string]AuthenticationSession
	m                      map[string]*sync.RWMutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		consentRequests:        map[string]ConsentRequest{},
		handledConsentRequests: map[string]HandledConsentRequest{},
		authRequests:           map[string]AuthenticationRequest{},
		handledAuthRequests:    map[string]HandledAuthenticationRequest{},
		authSessions:           map[string]AuthenticationSession{},
		m: map[string]*sync.RWMutex{
			"consentRequests":        new(sync.RWMutex),
			"handledConsentRequests": new(sync.RWMutex),
			"authRequests":           new(sync.RWMutex),
			"handledAuthRequests":    new(sync.RWMutex),
			"authSessions":           new(sync.RWMutex),
		},
	}
}

func (m *MemoryManager) CreateConsentRequest(c *ConsentRequest) error {
	m.m["consentRequests"].Lock()
	defer m.m["consentRequests"].Unlock()
	if _, ok := m.consentRequests[c.Challenge]; ok {
		return errors.New("Key already exists")
	}
	m.consentRequests[c.Challenge] = *c
	return nil
}

func (m *MemoryManager) GetConsentRequest(challenge string) (*ConsentRequest, error) {
	m.m["consentRequests"].RLock()
	defer m.m["consentRequests"].RUnlock()
	if c, ok := m.consentRequests[challenge]; ok {
		c.Client.ClientID = c.Client.ID
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) HandleConsentRequest(challenge string, r *HandledConsentRequest) (*ConsentRequest, error) {
	m.m["handledConsentRequests"].Lock()
	m.handledConsentRequests[r.Challenge] = *r
	m.m["handledConsentRequests"].Unlock()
	return m.GetConsentRequest(challenge)
}

func (m *MemoryManager) VerifyAndInvalidateConsentRequest(verifier string) (*HandledConsentRequest, error) {
	for _, c := range m.consentRequests {
		if c.Verifier == verifier {
			for _, h := range m.handledConsentRequests {
				if h.Challenge == c.Challenge {
					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Consent verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleConsentRequest(h.Challenge, &h); err != nil {
						return nil, err
					}

					c.Client.ClientID = c.Client.ID
					h.ConsentRequest = &c
					return &h, nil
				}
			}
		}
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) FindPreviouslyGrantedConsentRequests(client string, subject string) ([]HandledConsentRequest, error) {
	var rs []HandledConsentRequest
	for _, c := range m.handledConsentRequests {
		cr, err := m.GetConsentRequest(c.Challenge)
		if errors.Cause(err) == pkg.ErrNotFound {
			return nil, errors.WithStack(errNoPreviousConsentFound)
		} else if err != nil {
			return nil, err
		}

		if client != cr.Client.GetID() {
			continue
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

		cr.Client.ClientID = cr.Client.ID
		c.ConsentRequest = cr
		rs = append(rs, c)
	}
	if len(rs) == 0 {
		return []HandledConsentRequest{}, nil
	}

	return rs, nil
}

func (m *MemoryManager) GetAuthenticationSession(id string) (*AuthenticationSession, error) {
	m.m["authSessions"].RLock()
	defer m.m["authSessions"].RUnlock()
	if c, ok := m.authSessions[id]; ok {
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) CreateAuthenticationSession(a *AuthenticationSession) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	if _, ok := m.authSessions[a.ID]; ok {
		return errors.New("Key already exists")
	}
	m.authSessions[a.ID] = *a
	return nil
}

func (m *MemoryManager) DeleteAuthenticationSession(id string) error {
	m.m["authSessions"].Lock()
	defer m.m["authSessions"].Unlock()
	delete(m.authSessions, id)
	return nil
}

func (m *MemoryManager) CreateAuthenticationRequest(a *AuthenticationRequest) error {
	m.m["authRequests"].Lock()
	defer m.m["authRequests"].Unlock()
	if _, ok := m.authRequests[a.Challenge]; ok {
		return errors.New("Key already exists")
	}
	m.authRequests[a.Challenge] = *a
	return nil
}

func (m *MemoryManager) GetAuthenticationRequest(challenge string) (*AuthenticationRequest, error) {
	m.m["authRequests"].RLock()
	defer m.m["authRequests"].RUnlock()
	if c, ok := m.authRequests[challenge]; ok {
		c.Client.ClientID = c.Client.ID
		return &c, nil
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}

func (m *MemoryManager) HandleAuthenticationRequest(challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error) {
	m.m["handledAuthRequests"].Lock()
	m.handledAuthRequests[r.Challenge] = *r
	m.m["handledAuthRequests"].Unlock()
	return m.GetAuthenticationRequest(challenge)
}

func (m *MemoryManager) VerifyAndInvalidateAuthenticationRequest(verifier string) (*HandledAuthenticationRequest, error) {
	for _, c := range m.authRequests {
		if c.Verifier == verifier {
			for _, h := range m.handledAuthRequests {
				if h.Challenge == c.Challenge {
					if h.WasUsed {
						return nil, errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Authentication verifier has been used already"))
					}

					h.WasUsed = true
					if _, err := m.HandleAuthenticationRequest(h.Challenge, &h); err != nil {
						return nil, err
					}

					c.Client.ClientID = c.Client.ID
					h.AuthenticationRequest = &c
					return &h, nil
				}
			}
		}
	}
	return nil, errors.WithStack(pkg.ErrNotFound)
}
