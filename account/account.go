package account

type Account interface {
	GetID() string
	GetPassword() string
	GetUsername() string
	GetData() string
}

type DefaultAccount struct {
	ID       string `json:"id"`
	Username string `json:"username" valid:"required"`
	Password string `json:"-"`
	Data     string `json:"data,omitempty" valid:"json"`
}

func (a *DefaultAccount) GetID() string {
	return a.ID
}

func (a *DefaultAccount) GetPassword() string {
	return a.Password
}

func (a *DefaultAccount) GetUsername() string {
	return a.Username
}

func (a *DefaultAccount) GetData() string {
	return a.Data
}
