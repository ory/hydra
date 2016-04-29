package key

type AsymmetricKey struct {
	ID      string
	Public  []byte
	Private []byte
}

type SymmetricKey struct {
	ID  string
	Key []byte
}
