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
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2

import (
	"context"
	"sync"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/pkg/errors"
)

func NewFositeMemoryStore(m client.Manager, ls time.Duration) *FositeMemoryStore {
	return &FositeMemoryStore{
		AuthorizeCodes:      make(map[string]fosite.Requester),
		IDSessions:          make(map[string]fosite.Requester),
		AccessTokens:        make(map[string]fosite.Requester),
		RefreshTokens:       make(map[string]fosite.Requester),
		AccessTokenLifespan: ls,
		Manager:             m,
	}
}

type FositeMemoryStore struct {
	client.Manager

	AuthorizeCodes      map[string]fosite.Requester
	IDSessions          map[string]fosite.Requester
	AccessTokens        map[string]fosite.Requester
	RefreshTokens       map[string]fosite.Requester
	PKCES               map[string]fosite.Requester
	AccessTokenLifespan time.Duration

	sync.RWMutex
}

func (s *FositeMemoryStore) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.IDSessions[authorizeCode] = requester
	return nil
}

func (s *FositeMemoryStore) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	cl, ok := s.IDSessions[authorizeCode]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return cl, nil
}

func (s *FositeMemoryStore) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.IDSessions, authorizeCode)
	return nil
}

func (s *FositeMemoryStore) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.AuthorizeCodes[code] = req
	return nil
}

func (s *FositeMemoryStore) GetAuthorizeCodeSession(_ context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.AuthorizeCodes[code]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return rel, nil
}

func (s *FositeMemoryStore) DeleteAuthorizeCodeSession(_ context.Context, code string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.AuthorizeCodes, code)
	return nil
}

func (s *FositeMemoryStore) CreateAccessTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.AccessTokens[signature] = req
	return nil
}

func (s *FositeMemoryStore) GetAccessTokenSession(_ context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.AccessTokens[signature]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return rel, nil
}

func (s *FositeMemoryStore) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	s.Lock()
	defer s.Unlock()
	return s.deleteAccessTokenSession(ctx, signature)
}

func (s *FositeMemoryStore) deleteAccessTokenSession(_ context.Context, signature string) error {
	delete(s.AccessTokens, signature)
	return nil
}

func (s *FositeMemoryStore) CreateRefreshTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.RefreshTokens[signature] = req
	return nil
}

func (s *FositeMemoryStore) GetRefreshTokenSession(_ context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.RefreshTokens[signature]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return rel, nil
}

func (s *FositeMemoryStore) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	s.Lock()
	defer s.Unlock()
	return s.deleteRefreshTokenSession(ctx, signature)
}

func (s *FositeMemoryStore) deleteRefreshTokenSession(_ context.Context, signature string) error {
	delete(s.RefreshTokens, signature)
	return nil
}

func (s *FositeMemoryStore) CreateImplicitAccessTokenSession(ctx context.Context, code string, req fosite.Requester) error {
	return s.CreateAccessTokenSession(ctx, code, req)
}

func (s *FositeMemoryStore) RevokeRefreshToken(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()
	var found bool
	for sig, token := range s.RefreshTokens {
		if token.GetID() == id {
			if err := s.deleteRefreshTokenSession(ctx, sig); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		return errors.New("Not found")
	}
	return nil
}

func (s *FositeMemoryStore) RevokeAccessToken(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()
	var found bool
	for sig, token := range s.AccessTokens {
		if token.GetID() == id {
			if err := s.deleteAccessTokenSession(ctx, sig); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		return errors.New("Not found")
	}
	return nil
}

func (s *FositeMemoryStore) FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time) error {
	s.Lock()
	defer s.Unlock()

	now := time.Now()
	for sig, token := range s.AccessTokens {
		expiresAt := token.GetRequestedAt().Add(s.AccessTokenLifespan)
		isExpired := expiresAt.Before(now)
		isNotAfter := token.GetRequestedAt().Before(notAfter)

		if isExpired && isNotAfter {
			if err := s.deleteAccessTokenSession(ctx, sig); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *FositeMemoryStore) CreatePKCERequestSession(_ context.Context, code string, req fosite.Requester) error {
	s.PKCES[code] = req
	return nil
}

func (s *FositeMemoryStore) GetPKCERequestSession(_ context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.PKCES[code]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

func (s *FositeMemoryStore) DeletePKCERequestSession(_ context.Context, code string) error {
	delete(s.PKCES, code)
	return nil
}
