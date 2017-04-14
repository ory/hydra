package oauth2

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"github.com/go-redis/redis"
)

type FositeRedisStore struct {
	client.Manager
	DB        *redis.Client
	KeyPrefix string
}

type redisSchema struct {
	ID            string           `json:"id"`
	RequestedAt   time.Time        `json:"requestedAt"`
	Client        *client.Client   `json:"client"`
	Scopes        fosite.Arguments `json:"scopes"`
	GrantedScopes fosite.Arguments `json:"grantedScopes"`
	Form          url.Values       `json:"form"`
	Session       json.RawMessage  `json:"session"`
}

func (s *FositeRedisStore) redisKey(fields ...string) string {
	return s.KeyPrefix + strings.Join(fields, ":")
}

const (
	redisHydraKey = "hydra:oauth2"
	redisOpenID   = "oidc"
	redisAccess   = "access"
	redisRefresh  = "refresh"
	redisCode     = "code"
)

func redisCreateTokenSession(pipe *redis.Pipeline, req fosite.Requester, key, setKey, signature string) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return errors.WithStack(err)
	}

	pipe.HSet(key, signature, string(payload))
	pipe.SAdd(setKey, signature)
	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return nil
}

func (s *FositeRedisStore) hGet(hash, key string, proto fosite.Session) (*fosite.Request, error) {
	resp, err := s.DB.HGet(hash, key).Bytes()
	if err != nil {
		return nil, err
	}

	var schema redisSchema
	if err := json.Unmarshal(resp, &schema); err != nil {
		return nil, err
	}

	if proto != nil {
		if err := json.Unmarshal(schema.Session, proto); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &fosite.Request{
		ID:            schema.ID,
		RequestedAt:   schema.RequestedAt,
		Client:        schema.Client,
		Scopes:        schema.Scopes,
		GrantedScopes: schema.GrantedScopes,
		Form:          schema.Form,
		Session:       proto,
	}, nil
}

func (s *FositeRedisStore) hSet(hash, key string, requester fosite.Requester) error {
	payload, err := json.Marshal(requester)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := s.DB.HSet(hash, key, string(payload)).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *FositeRedisStore) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, req fosite.Requester) error {
	return s.hSet(s.redisKey(redisOpenID), authorizeCode, req)
}

func (s *FositeRedisStore) GetOpenIDConnectSession(_ context.Context, authorizeCode string, req fosite.Requester) (fosite.Requester, error) {
	session, err := s.hGet(s.redisKey(redisOpenID), authorizeCode, req.GetSession())
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return session, nil
}

func (s *FositeRedisStore) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	return s.DB.HDel(s.redisKey(redisOpenID), authorizeCode).Err()
}

func (s *FositeRedisStore) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	return s.hSet(s.redisKey(redisCode), code, req)
}

func (s *FositeRedisStore) GetAuthorizeCodeSession(_ context.Context, code string, sess fosite.Session) (fosite.Requester, error) {
	session, err := s.hGet(s.redisKey(redisCode), code, sess)
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return session, nil
}

func (s *FositeRedisStore) DeleteAuthorizeCodeSession(_ context.Context, code string) error {
	return s.DB.HDel(s.redisKey(redisCode), code).Err()
}

func (s *FositeRedisStore) CreateAccessTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	pipe := s.DB.Pipeline()
	defer pipe.Close()

	return redisCreateTokenSession(
		pipe,
		req, s.redisKey(redisAccess),
		s.redisKey(redisAccess, req.GetID()),
		signature,
	)
}

func (s *FositeRedisStore) GetAccessTokenSession(_ context.Context, signature string, sess fosite.Session) (fosite.Requester, error) {
	session, err := s.hGet(s.redisKey(redisAccess), signature, sess)
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return session, nil
}

func (s *FositeRedisStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	sess, err := s.hGet(s.redisKey(redisAccess), signature, nil)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return errors.WithStack(err)
	}

	pipe := s.DB.Pipeline()
	defer pipe.Close()

	pipe.HDel(s.redisKey(redisAccess), signature)
	pipe.SRem(s.redisKey(redisAccess, sess.GetID()), signature)
	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return nil
}

func (s *FositeRedisStore) CreateRefreshTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	pipe := s.DB.Pipeline()
	defer pipe.Close()

	return redisCreateTokenSession(
		pipe,
		req,
		s.redisKey(redisRefresh),
		s.redisKey(redisRefresh, req.GetID()),
		signature,
	)
}

func (s *FositeRedisStore) GetRefreshTokenSession(_ context.Context, signature string, sess fosite.Session) (fosite.Requester, error) {
	return s.hGet(s.redisKey(redisRefresh), signature, sess)
}

func (s *FositeRedisStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	sess, err := s.hGet(s.redisKey(redisRefresh), signature, nil)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return errors.WithStack(err)
	}

	pipe := s.DB.Pipeline()
	defer pipe.Close()

	pipe.HDel(s.redisKey(redisRefresh), signature)
	pipe.SRem(s.redisKey(redisRefresh, sess.GetID()), signature)

	_, err = pipe.Exec()
	return err
}

func (s *FositeRedisStore) CreateImplicitAccessTokenSession(ctx context.Context, code string, req fosite.Requester) error {
	return s.CreateAccessTokenSession(ctx, code, req)
}

func (s *FositeRedisStore) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, req fosite.Requester) error {
	if err := s.DeleteAuthorizeCodeSession(ctx, authorizeCode); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, req); err != nil {
		return err
	}

	if refreshSignature == "" {
		return nil
	}

	if err := s.CreateRefreshTokenSession(ctx, refreshSignature, req); err != nil {
		return err
	}

	return nil
}

func (s *FositeRedisStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, req fosite.Requester) error {
	if err := s.DeleteRefreshTokenSession(ctx, originalRefreshSignature); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, req); err != nil {
		return err
	} else if err := s.CreateRefreshTokenSession(ctx, refreshSignature, req); err != nil {
		return err
	}

	return nil
}

func (s *FositeRedisStore) RevokeRefreshToken(ctx context.Context, id string) error {
	pipe := s.DB.Pipeline()
	defer pipe.Close()

	refreshSet := s.redisKey(redisRefresh, id)
	iter := s.DB.SScan(refreshSet, 0, "", 0).Iterator()
	for iter.Next() {
		sig := iter.Val()
		pipe.HDel(redisRefresh, sig)
		pipe.SRem(refreshSet, sig)
	}
	if err := iter.Err(); err != nil {
		return err
	}

	_, err := pipe.Exec()
	return err
}

func (s *FositeRedisStore) RevokeAccessToken(ctx context.Context, id string) error {
	pipe := s.DB.Pipeline()
	defer pipe.Close()

	accessSet := s.redisKey(redisAccess, id)
	iter := s.DB.SScan(accessSet, 0, "", 0).Iterator()
	for iter.Next() {
		sig := iter.Val()
		pipe.HDel(redisAccess, sig)
		pipe.SRem(accessSet, sig)
	}
	if err := iter.Err(); err != nil {
		return err
	}

	_, err := pipe.Exec()
	return err
}
