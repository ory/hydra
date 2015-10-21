package account

type Account interface {
	GetID() string
	GetPassword() string
	GetEmail() string
	GetData() string
}
