package account

type Storage interface {
	Create(id, email, password string, data string) (Account, error)

	Get(id string) (Account, error)

	Delete(id string) error

	UpdatePassword(id, oldPassword, newPassword string) (Account, error)

	UpdateEmail(id, password, email string) (Account, error)

	UpdateData(id, data string) (Account, error)

	Authenticate(email, password string) (Account, error)

	FindByProvider(provider, subject string) (Account, error)
}
