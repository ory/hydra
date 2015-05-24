package hash

type Hasher interface {
	Compare(hash, data string) error
	Hash(data string) (string, error)
}
