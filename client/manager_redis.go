package client

import (
	"encoding/json"

	"github.com/imdario/mergo"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/go-redis/redis"
)

type RedisManager struct {
	DB        *redis.Client
	Hasher    fosite.Hasher
	KeyPrefix string
}

const redisClientTemplate = "hydra:client"

func (m *RedisManager) redisClientKey() string {
	return m.KeyPrefix + redisClientTemplate
}

func (m *RedisManager) GetConcreteClient(id string) (*Client, error) {
	resp, err := m.DB.HGet(m.redisClientKey(), id).Bytes()
	if err == redis.Nil {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	var d Client
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, errors.WithStack(err)
	}

	return &d, nil
}

func (m *RedisManager) GetClient(id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *RedisManager) UpdateClient(c *Client) error {
	o, err := m.GetClient(c.ID)
	if err != nil {
		return err
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(c.Secret))
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = string(h)
	}
	if err := mergo.Merge(c, o); err != nil {
		return errors.WithStack(err)
	}

	s, err := json.Marshal(c)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := m.DB.HSet(m.redisClientKey(), c.ID, string(s)).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *RedisManager) Authenticate(id string, secret []byte) (*Client, error) {
	c, err := m.GetConcreteClient(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (m *RedisManager) CreateClient(c *Client) error {
	if c.ID == "" {
		c.ID = uuid.New()
	}

	hash, err := m.Hasher.Hash([]byte(c.Secret))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = string(hash)

	s, err := json.Marshal(c)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := m.DB.HSetNX(m.redisClientKey(), c.ID, string(s)).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *RedisManager) DeleteClient(id string) error {
	if _, err := m.DB.HDel(m.redisClientKey(), id).Result(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *RedisManager) GetClients() (map[string]Client, error) {
	clients := make(map[string]Client)

	iter := m.DB.HScan(m.redisClientKey(), 0, "", 0).Iterator()
	for iter.Next() {
		if !iter.Next() {
			break
		}

		resp := iter.Val()

		var d Client
		if err := json.Unmarshal([]byte(resp), &d); err != nil {
			return nil, errors.WithStack(err)
		}

		clients[d.ID] = d
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
