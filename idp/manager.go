package idp

type IdentityProvider struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	Healthy  bool   `json:"healthy"`
}

type Manager interface {
	Register(*IdentityProvider) error

	All() (map[string]*IdentityProvider, error)

	Get(id string) (*IdentityProvider, error)

	Deregister(id string) error
}
