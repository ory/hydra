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

package oauth2

import (
	"sync"

	"context"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/pkg/errors"
)

type FositeMemoryStore struct {
	client.Manager

	AuthorizeCodes map[string]fosite.Requester
	IDSessions     map[string]fosite.Requester
	AccessTokens   map[string]fosite.Requester
	RefreshTokens  map[string]fosite.Requester

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
