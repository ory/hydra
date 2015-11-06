package account

import "github.com/go-errors/errors"

var ErrNotFound = errors.New("Not found")

type Account interface {
	GetID() string
	GetPassword() string
	GetEmail() string
	GetData() string
}

type DefaultAccount struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Data     string `json:"data"`
}

func (a *DefaultAccount) GetID() string {
	return a.ID
}

func (a *DefaultAccount) GetPassword() string {
	return a.Password
}

func (a *DefaultAccount) GetEmail() string {
	return a.Email
}

func (a *DefaultAccount) GetData() string {
	return a.Data
}
