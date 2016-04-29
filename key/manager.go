package key

type Manager interface {
	CreateAsymmetricKey(id string) (*AsymmetricKey, error)

	GetAsymmetricKey(id string) (*AsymmetricKey, error)

	DeleteAsymmetricKey(id string) error

	CreateSymmetricKey(id string) (*SymmetricKey, error)

	GetSymmetricKey(id string) (*SymmetricKey, error)

	DeleteSymmetricKey(id string) error
}
