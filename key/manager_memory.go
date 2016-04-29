package key

import "github.com/go-errors/errors"

type MemoryManager struct {
	Strategy KeyStrategy

	AsymmetricKeys map[string]*AsymmetricKey

	SymmetricKeys map[string]*SymmetricKey
}

func (m *MemoryManager) CreateAsymmetricKey(id string) (*AsymmetricKey, error) {
	key, err := m.Strategy.AsymmetricKey(id)
	if err != nil {
		return nil, err
	}

	m.AsymmetricKeys[id] = key
	return key, nil
}

func (m *MemoryManager) DeleteAsymmetricKey(id string) error {
	delete(m.AsymmetricKeys, id)
	return nil
}

func (m *MemoryManager) GetAsymmetricKey(id string) (*AsymmetricKey, error) {
	key, ok := m.AsymmetricKeys[id]
	if !ok {
		return nil, errors.New("Key not found")
	}

	return key, nil
}

func (m *MemoryManager) CreateSymmetricKey(id string) (*SymmetricKey, error) {
	key, err := m.Strategy.SymmetricKey(id)
	if err != nil {
		return nil, err
	}

	m.SymmetricKeys[id] = key
	return key, nil

}

func (m *MemoryManager) DeleteSymmetricKey(id string) error {
	delete(m.SymmetricKeys, id)
	return nil
}

func (m *MemoryManager) GetSymmetricKey(id string) (*SymmetricKey, error) {
	key, ok := m.SymmetricKeys[id]
	if !ok {
		return nil, errors.New("Key not found")
	}

	return key, nil
}
