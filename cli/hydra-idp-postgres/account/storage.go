package account

// UpdatePasswordRequest is used to pass the request data to Storage.UpdatePassword.
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" valid:"required"`
	NewPassword     string `json:"newPassword" valid:"required"`
}

// UpdateUsernameRequest is used to pass the request data to Storage.UpdateUsername.
type UpdateUsernameRequest struct {
	Password string `json:"password" valid:"required"`
	Username string `json:"username" valid:"required"`
}

// UpdateDataRequest is used to pass an accounts data to Storage.UpdateData.
type UpdateDataRequest struct {
	Data string `json:"data" valid:"json,required"`
}

// CreateAccountRequest is used to pass an accounts data to Storage.Create.
type CreateAccountRequest struct {
	ID       string `json:"id" valid:"uuid4"`
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
	Data     string `json:"data" valid:"json"`
}

// Storage manages accounts in a persistent fashion.
type Storage interface {

	// Create creates a new account
	Create(r CreateAccountRequest) (Account, error)

	// Get fetches an account by its ID.
	Get(id string) (Account, error)

	// Delete deletes an account by its ID.
	Delete(id string) error

	// UpdatePassword updates an account's password.
	UpdatePassword(id string, r UpdatePasswordRequest) (Account, error)

	// UpdateUsername updates an account's username.
	UpdateUsername(id string, r UpdateUsernameRequest) (Account, error)

	// UpdateData updates an account's extra data (e.g. profile picture)
	UpdateData(id string, r UpdateDataRequest) (Account, error)

	// Authenticate fetches an account by its username and password
	Authenticate(username, password string) (Account, error)
}
