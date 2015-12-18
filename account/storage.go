package account

type Storage interface {
	Create(id, username, password string, data string) (Account, error)

	Get(id string) (Account, error)

	Delete(id string) error

	UpdatePassword(id, oldPassword, newPassword string) (Account, error)

	UpdateUsername(id, password, username string) (Account, error)

	UpdateData(id, data string) (Account, error)

	Authenticate(username, password string) (Account, error)
}
