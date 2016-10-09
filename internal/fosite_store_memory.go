package internal

import (
	"sync"

	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type FositeMemoryStore struct {
	client.Manager

	AuthorizeCodes map[string]fosite.Requester
	IDSessions     map[string]fosite.Requester
	AccessTokens   map[string]fosite.Requester
	Implicit       map[string]fosite.Requester
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

func (s *FositeMemoryStore) GetAuthorizeCodeSession(_ context.Context, code string, _ interface{}) (fosite.Requester, error) {
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

func (s *FositeMemoryStore) GetAccessTokenSession(_ context.Context, signature string, _ interface{}) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.AccessTokens[signature]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return rel, nil
}

func (s *FositeMemoryStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.AccessTokens, signature)
	return nil
}

func (s *FositeMemoryStore) CreateRefreshTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.RefreshTokens[signature] = req
	return nil
}

func (s *FositeMemoryStore) GetRefreshTokenSession(_ context.Context, signature string, _ interface{}) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.RefreshTokens[signature]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return rel, nil
}

func (s *FositeMemoryStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.RefreshTokens, signature)
	return nil
}

func (s *FositeMemoryStore) CreateImplicitAccessTokenSession(_ context.Context, code string, req fosite.Requester) error {
	s.Lock()
	defer s.Unlock()
	s.Implicit[code] = req
	return nil
}

func (s *FositeMemoryStore) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := s.DeleteAuthorizeCodeSession(ctx, authorizeCode); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	}

	if refreshSignature == "" {
		return nil
	}

	if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (s *FositeMemoryStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := s.DeleteRefreshTokenSession(ctx, originalRefreshSignature); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}
