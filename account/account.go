package account

type Account interface {
	GetID() string

	GetPassword() string
	SetPassword(password string)

	GetEmail() string
	SetEmail(email string)
}
