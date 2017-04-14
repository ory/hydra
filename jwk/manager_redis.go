package jwk

import (
	"encoding/json"
	"fmt"

	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"github.com/go-redis/redis"
)

type RedisManager struct {
	DB        *redis.Client
	Cipher    *AEAD
	KeyPrefix string
}

func (m *RedisManager) redisJWKKey(set string) string {
	return m.KeyPrefix + fmt.Sprintf("hydra:jwk:%s", set)
}

func (m *RedisManager) addKey(set string, key *jose.JsonWebKey, pipe *redis.Pipeline) error {
	payload, err := json.Marshal(key)
	if err != nil {
		return err
	}

	encrypted, err := m.Cipher.Encrypt(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	pipe.HSet(m.redisJWKKey(set), key.KeyID, encrypted)

	return nil
}

func (m *RedisManager) AddKey(set string, key *jose.JsonWebKey) error {
	pipe := m.DB.Pipeline()
	defer pipe.Close()

	if err := m.addKey(set, key, pipe); err != nil {
		return errors.WithStack(err)
	}

	if _, err := pipe.Exec(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *RedisManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	pipe := m.DB.Pipeline()
	defer pipe.Close()

	for _, key := range keys.Keys {
		if err := m.addKey(set, &key, pipe); err != nil {
			return errors.WithStack(err)
		}
	}

	if _, err := pipe.Exec(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *RedisManager) getKeySet(set string) (*jose.JsonWebKeySet, error) {
	var keySet jose.JsonWebKeySet

	iter := m.DB.HScan(m.redisJWKKey(set), 0, "", 0).Iterator()
	for iter.Next() {
		if !iter.Next() {
			break
		}
		encryptedJWK := iter.Val()

		jwk, err := m.Cipher.Decrypt(encryptedJWK)
		if err != nil {
			return nil, err
		}

		var key jose.JsonWebKey
		if err := json.Unmarshal(jwk, &key); err != nil {
			return nil, err
		}

		keySet.Keys = append(keySet.Keys, key)
	}

	if len(keySet.Keys) == 0 {
		return nil, pkg.ErrNotFound
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &keySet, nil
}

func (m *RedisManager) GetKey(set, kid string) (*jose.JsonWebKeySet, error) {
	encryptedJWK, err := m.DB.HGet(m.redisJWKKey(set), kid).Result()
	if err == redis.Nil {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	jwk, err := m.Cipher.Decrypt(encryptedJWK)
	if err != nil {
		return nil, err
	}

	var key jose.JsonWebKey
	if err := json.Unmarshal(jwk, &key); err != nil {
		return nil, errors.WithStack(err)
	}

	return &jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{key},
	}, nil
}

func (m *RedisManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	keys, err := m.getKeySet(set)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return keys, nil
}

func (m *RedisManager) DeleteKey(set, kid string) error {
	return m.DB.HDel(m.redisJWKKey(set), kid).Err()
}

func (m *RedisManager) DeleteKeySet(set string) error {
	return m.DB.Del(m.redisJWKKey(set)).Err()
}
