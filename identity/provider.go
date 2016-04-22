package identity

type IdentityProvider interface {
	GetID() string

	Authenticate(id, password string) (*Identity, error)

	GetIdentity(id string) (*Identity, error)
}
