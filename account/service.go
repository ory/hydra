package account

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/ory-am/go-iam/hash"
)

const passwordWorkFactor = 12

type Service struct {
	Storage Storage
	Hash    hash.Hasher
}

func (s *Service) Create(email, password string) (Account, error) {
	if hash, err := s.Hash.Hash(password); err != nil {
		return nil, err
	} else {
		return s.Storage.Create(uuid.NewRandom().String(), email, hash)
	}
}

func (s *Service) Get(id string) (Account, error) {
	return s.Storage.Get(id)
}

func (s *Service) UpdatePassword(id, old, new string) (Account, error) {
	if acc, err := s.Storage.Get(id); err != nil {
		return nil, err
	} else {
		if err := s.Hash.Compare(acc.GetPassword(), old); err != nil {
			return nil, err
		}
		if hashed, err := s.Hash.Hash(new); err != nil {
			return nil, err
		} else {
			acc.SetPassword(hashed)
			s.Storage.Update(acc)
			return acc, err
		}
	}
}

func (s *Service) UpdateEmail(id, email string) (Account, error) {
	if acc, err := s.Storage.Get(id); err != nil {
		return nil, err
	} else {
		acc.SetEmail(email)
		s.Storage.Update(acc)
		return acc, err
	}
}
