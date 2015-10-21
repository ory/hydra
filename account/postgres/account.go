package postgres

type Account struct {
	ID       string
	Email    string
	Password string
	Data     string
}

func (a *Account) GetID() string {
	return a.ID
}

func (a *Account) GetPassword() string {
	return a.Password
}

func (a *Account) GetEmail() string {
	return a.Email
}

func (a *Account) GetData() string {
	return a.Data
}
