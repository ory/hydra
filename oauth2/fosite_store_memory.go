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

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/x/sqlcon"
)

func NewFositeMemoryStore(m client.Manager, ls time.Duration) *FositeMemoryStore {
	return &FositeMemoryStore{
		AuthorizeCodes:      make(map[string]authorizeCode),
		IDSessions:          make(map[string]fosite.Requester),
		AccessTokens:        make(map[string]fosite.Requester),
		PKCES:               make(map[string]fosite.Requester),
		RefreshTokens:       make(map[string]fosite.Requester),
		AccessTokenLifespan: ls,
		Manager:             m,
	}
}

type FositeMemoryStore struct {
	client.Manager

	AuthorizeCodes      map[string]authorizeCode
	IDSessions          map[string]fosite.Requester
	AccessTokens        map[string]fosite.Requester
	RefreshTokens       map[string]fosite.Requester
	PKCES               map[string]fosite.Requester
	AccessTokenLifespan time.Duration

	sync.RWMutex
}

type authorizeCode struct {
	active bool
	fosite.Requester
}

func (s *FositeMemoryStore) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.IDSessions[authorizeCode] = requester
	return nil
}

func (s *FositeMemoryStore) GetOpenIDConnectSession(ctx context.Context, code string, requester fosite.Requester) (fosite.Requester, error) {
	s.RLock()
	rel, ok := s.IDSessions[code]
	s.RUnlock()

	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	if _, err := s.GetClient(ctx, rel.GetClient().GetID()); errors.Cause(err) == sqlcon.ErrNoRows {
		s.Lock()
		delete(s.IDSessions, code)
		s.Unlock()
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return rel, nil
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
	s.AuthorizeCodes[code] = authorizeCode{active: true, Requester: req}
	return nil
}

func (s *FositeMemoryStore) GetAuthorizeCodeSession(ctx context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	rel, ok := s.AuthorizeCodes[code]
	s.RUnlock()

	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	if !rel.active {
		return rel.Requester, errors.WithStack(fosite.ErrInvalidatedAuthorizeCode)
	}

	if _, err := s.GetClient(ctx, rel.GetClient().GetID()); errors.Cause(err) == sqlcon.ErrNoRows {
		s.Lock()
		delete(s.AuthorizeCodes, code)
		s.Unlock()
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return rel.Requester, nil
}

func (s *FositeMemoryStore) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	s.Lock()
	defer s.Unlock()

	rel, ok := s.AuthorizeCodes[code]
	if !ok {
		return fosite.ErrNotFound
	}
	rel.active = false
	s.AuthorizeCodes[code] = rel
	return nil
}

func (s *FositeMemoryStore) CreateAccessTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.AccessTokens[signature] = req
	return nil
}

func (s *FositeMemoryStore) GetAccessTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	rel, ok := s.AccessTokens[signature]
	s.RUnlock()

	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	if _, err := s.GetClient(ctx, rel.GetClient().GetID()); errors.Cause(err) == sqlcon.ErrNoRows {
		s.Lock()
		delete(s.AccessTokens, signature)
		s.Unlock()
		return nil, err
	} else if err != nil {
		return nil, err
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

func (s *FositeMemoryStore) GetRefreshTokenSession(ctx context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	rel, ok := s.RefreshTokens[signature]
	s.RUnlock()

	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	if _, err := s.GetClient(ctx, rel.GetClient().GetID()); errors.Cause(err) == sqlcon.ErrNoRows {
		s.Lock()
		delete(s.RefreshTokens, signature)
		s.Unlock()
		return nil, err
	} else if err != nil {
		return nil, err
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
		return errors.WithStack(fosite.ErrNotFound)
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
		return errors.WithStack(fosite.ErrNotFound)
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
	s.Lock()
	s.PKCES[code] = req
	s.Unlock()
	return nil
}

func (s *FositeMemoryStore) GetPKCERequestSession(ctx context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	s.RLock()
	rel, ok := s.PKCES[code]
	s.RUnlock()
	if !ok {
		return nil, fosite.ErrNotFound
	}

	if _, err := s.GetClient(ctx, rel.GetClient().GetID()); errors.Cause(err) == sqlcon.ErrNoRows {
		s.Lock()
		delete(s.RefreshTokens, code)
		s.Unlock()
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return rel, nil
}

func (s *FositeMemoryStore) DeletePKCERequestSession(_ context.Context, code string) error {
	s.Lock()
	delete(s.PKCES, code)
	s.Unlock()
	return nil
}
