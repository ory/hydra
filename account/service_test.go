package account

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/ory-am/go-iam/hash"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	s := Service{
		Storage: new(storageMock),
		Hash:    new(hash.BCrypt),
	}
	acc, err := s.Create("foo", "bar")
	assert.Nil(t, err)
	assert.NotEmpty(t, acc.GetID())
	assert.Equal(t, acc.GetEmail(), "foo")
	assert.NotEqual(t, acc.GetPassword(), "bar")
	assert.NotNil(t, acc.GetPassword())
}

func TestGet(t *testing.T) {
	s := Service{
		Storage: new(storageMock),
		Hash:    new(hash.BCrypt),
	}
	acc, err := s.Get("foo")
	assert.Nil(t, err)
	assert.NotNil(t, acc)
}

func TestUpdatePassword(t *testing.T) {
	s := Service{
		Storage: new(storageMock),
		Hash:    new(hash.BCrypt),
	}
	acc, err := s.Get("baz")
	assert.Nil(t, err)
	oldHash := acc.GetPassword()
	acc, err = s.UpdatePassword("foo", "bar", "baz")
	assert.Nil(t, err)
	assert.NotEqual(t, acc.GetPassword(), oldHash)
}

type storageMock struct{}

func (s *storageMock) Create(id, email, password string) (Account, error) {
	return &accountMock{
		id:       uuid.NewRandom().String(),
		email:    email,
		password: password,
	}, nil
}

func (s *storageMock) Update(account Account) error {
	return nil
}

func (s *storageMock) Get(id string) (Account, error) {
	b := new(hash.BCrypt)
	h, _ := b.Hash("bar")
	return &accountMock{
		id:       id,
		email:    "foo@bar",
		password: h,
	}, nil
}

type accountMock struct {
	id       string
	email    string
	password string
}

func (a *accountMock) GetEmail() string {
	return a.email
}
func (a *accountMock) GetID() string {
	return a.id
}
func (a *accountMock) GetPassword() string {
	return a.password
}

func (a *accountMock) SetEmail(email string) {
	a.email = email
}
func (a *accountMock) SetPassword(password string) {
	a.password = password
}
