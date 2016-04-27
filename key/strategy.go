package key

type AsymmetricKeyStrategy interface {
	AsymmetricKey(id string) (*AsymmetricKey, error)
}

type SymmetricKeyStrategy interface {
	SymmetricKey(id string) (*SymmetricKey, error)
}

type KeyStrategy interface {
	AsymmetricKeyStrategy
	SymmetricKeyStrategy
}

type DefaultKeyStrategy struct {
	AsymmetricKeyStrategy
	SymmetricKeyStrategy
}