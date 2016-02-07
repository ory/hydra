package account

type Account interface {
	GetID() string
	GetPassword() string
	GetUsername() string
	GetData() string
}

type DefaultAccount struct {
	ID       string `json:"id" gorethink:"id"`
	Username string `json:"username" valid:"required" gorethink:"username"`
	Password string `json:"-" gorethink:"password"`
	Data     string `json:"data,omitempty" valid:"json" gorethink:"data"`
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
