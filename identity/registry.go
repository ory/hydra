package identity

type IdentityProviderRegistry interface {
	AddProvider(provider IdentityProvider)

	RemoveProvider(provider IdentityProvider)

	Authenticate(id, password, otp string) (*Identity, error)

	GetIdentity(id string) (*Identity, error)
}
