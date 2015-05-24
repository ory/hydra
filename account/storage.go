package account

type Storage interface {
	Create(id, email, password string) (Account, error)
	Update(account Account) error
	Get(id string) (Account, error)
}
