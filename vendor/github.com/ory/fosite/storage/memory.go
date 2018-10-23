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
 *
 */

package storage

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

type MemoryUserRelation struct {
	Username string
	Password string
}

type MemoryStore struct {
	Clients        map[string]fosite.Client
	AuthorizeCodes map[string]StoreAuthorizeCode
	IDSessions     map[string]fosite.Requester
	AccessTokens   map[string]fosite.Requester
	Implicit       map[string]fosite.Requester
	RefreshTokens  map[string]fosite.Requester
	PKCES          map[string]fosite.Requester
	Users          map[string]MemoryUserRelation
	// In-memory request ID to token signatures
	AccessTokenRequestIDs  map[string]string
	RefreshTokenRequestIDs map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Clients:                make(map[string]fosite.Client),
		AuthorizeCodes:         make(map[string]StoreAuthorizeCode),
		IDSessions:             make(map[string]fosite.Requester),
		AccessTokens:           make(map[string]fosite.Requester),
		Implicit:               make(map[string]fosite.Requester),
		RefreshTokens:          make(map[string]fosite.Requester),
		PKCES:                  make(map[string]fosite.Requester),
		Users:                  make(map[string]MemoryUserRelation),
		AccessTokenRequestIDs:  make(map[string]string),
		RefreshTokenRequestIDs: make(map[string]string),
	}
}

type StoreAuthorizeCode struct {
	active bool
	fosite.Requester
}

func NewExampleStore() *MemoryStore {
	return &MemoryStore{
		IDSessions: make(map[string]fosite.Requester),
		Clients: map[string]fosite.Client{
			"my-client": &fosite.DefaultClient{
				ID:            "my-client",
				Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
				RedirectURIs:  []string{"http://localhost:3846/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scopes:        []string{"fosite", "openid", "photos", "offline"},
			},
			"encoded:client": &fosite.DefaultClient{
				ID:            "encoded:client",
				Secret:        []byte(`$2a$10$A7M8b65dSSKGHF0H2sNkn.9Z0hT8U1Nv6OWPV3teUUaczXkVkxuDS`), // = "encoded&password"
				RedirectURIs:  []string{"http://localhost:3846/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scopes:        []string{"fosite", "openid", "photos", "offline"},
			},
		},
		Users: map[string]MemoryUserRelation{
			"peter": {
				// This store simply checks for equality, a real storage implementation would obviously use
				// a hashing algorithm for encrypting the user password.
				Username: "peter",
				Password: "secret",
			},
		},
		AuthorizeCodes:         map[string]StoreAuthorizeCode{},
		Implicit:               map[string]fosite.Requester{},
		AccessTokens:           map[string]fosite.Requester{},
		RefreshTokens:          map[string]fosite.Requester{},
		PKCES:                  map[string]fosite.Requester{},
		AccessTokenRequestIDs:  map[string]string{},
		RefreshTokenRequestIDs: map[string]string{},
	}
}

func (s *MemoryStore) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	s.IDSessions[authorizeCode] = requester
	return nil
}

func (s *MemoryStore) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	cl, ok := s.IDSessions[authorizeCode]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return cl, nil
}

func (s *MemoryStore) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	delete(s.IDSessions, authorizeCode)
	return nil
}

func (s *MemoryStore) GetClient(_ context.Context, id string) (fosite.Client, error) {
	cl, ok := s.Clients[id]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return cl, nil
}

func (s *MemoryStore) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	s.AuthorizeCodes[code] = StoreAuthorizeCode{active: true, Requester: req}
	return nil
}

func (s *MemoryStore) GetAuthorizeCodeSession(_ context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.AuthorizeCodes[code]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	if !rel.active {
		return rel, fosite.ErrInvalidatedAuthorizeCode
	}

	return rel.Requester, nil
}

func (s *MemoryStore) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	rel, ok := s.AuthorizeCodes[code]
	if !ok {
		return fosite.ErrNotFound
	}
	rel.active = false
	s.AuthorizeCodes[code] = rel
	return nil
}

func (s *MemoryStore) DeleteAuthorizeCodeSession(_ context.Context, code string) error {
	delete(s.AuthorizeCodes, code)
	return nil
}

func (s *MemoryStore) CreatePKCERequestSession(_ context.Context, code string, req fosite.Requester) error {
	s.PKCES[code] = req
	return nil
}

func (s *MemoryStore) GetPKCERequestSession(_ context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.PKCES[code]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

func (s *MemoryStore) DeletePKCERequestSession(_ context.Context, code string) error {
	delete(s.PKCES, code)
	return nil
}

func (s *MemoryStore) CreateAccessTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.AccessTokens[signature] = req
	s.AccessTokenRequestIDs[req.GetID()] = signature
	return nil
}

func (s *MemoryStore) GetAccessTokenSession(_ context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.AccessTokens[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

func (s *MemoryStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	delete(s.AccessTokens, signature)
	return nil
}

func (s *MemoryStore) CreateRefreshTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.RefreshTokens[signature] = req
	s.RefreshTokenRequestIDs[req.GetID()] = signature
	return nil
}

func (s *MemoryStore) GetRefreshTokenSession(_ context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.RefreshTokens[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

func (s *MemoryStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	delete(s.RefreshTokens, signature)
	return nil
}

func (s *MemoryStore) CreateImplicitAccessTokenSession(_ context.Context, code string, req fosite.Requester) error {
	s.Implicit[code] = req
	return nil
}

func (s *MemoryStore) Authenticate(_ context.Context, name string, secret string) error {
	rel, ok := s.Users[name]
	if !ok {
		return fosite.ErrNotFound
	}
	if rel.Password != secret {
		return errors.New("Invalid credentials")
	}
	return nil
}

func (s *MemoryStore) RevokeRefreshToken(ctx context.Context, requestID string) error {
	if signature, exists := s.RefreshTokenRequestIDs[requestID]; exists {
		s.DeleteRefreshTokenSession(ctx, signature)
		s.DeleteAccessTokenSession(ctx, signature)
	}
	return nil
}

func (s *MemoryStore) RevokeAccessToken(ctx context.Context, requestID string) error {
	if signature, exists := s.AccessTokenRequestIDs[requestID]; exists {
		s.DeleteAccessTokenSession(ctx, signature)
	}
	return nil
}
