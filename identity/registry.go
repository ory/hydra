package identity

type IdentityProviderRegistry interface {
	AddProvider(provider IdentityProvider)

	RemoveProvider(provider IdentityProvider)

	Authenticate(id, password string) (Identity, error)

	IsIdentityAuthenticable(id string) error
}
