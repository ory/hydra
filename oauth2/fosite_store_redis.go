package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ory/hydra/v2/client"
	"golang.org/x/text/language"
	"gopkg.in/square/go-jose.v2"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type FositeRedisStore struct {
	//client.Manager
	//fosite.CoreStorage
	DB        redis.UniversalClient
	KeyPrefix string
}

const (
	prefixOIDC         = "oidc"
	prefixAccess       = "access"
	prefixRefresh      = "refresh"
	prefixCode         = "code"
	prefixPKCE         = "pkce"
	prefixClient       = "client"
	prefixJTIBlocklist = "block-jti"
)

type redisSchema struct {
	ID                string           `json:"id"`
	RequestedAt       time.Time        `json:"requestedAt"`
	Client            *client.Client   `json:"client"`
	RequestedScope    fosite.Arguments `json:"scopes"`
	GrantedScope      fosite.Arguments `json:"grantedScopes"`
	Form              url.Values       `json:"form"`
	Session           json.RawMessage  `json:"session"`
	RequestedAudience fosite.Arguments `json:"requestedAudience"`
	GrantedAudience   fosite.Arguments `json:"grantedAudience"`
	Lang              language.Tag     `json:"-"`

	// field to track revocation
	Active bool `json:"active"`
}

func (s FositeRedisStore) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, req fosite.Requester) error {
	return s.setRequest(ctx, s.redisKey(prefixOIDC), authorizeCode, req)
}

func (s FositeRedisStore) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, req fosite.Requester) (fosite.Requester, error) {
	session, _, err := s.getRequest(ctx, s.redisKey(prefixOIDC), authorizeCode, req.GetSession())
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return session, nil
}

func (s FositeRedisStore) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	fullKey := s.redisKey(prefixClient, id)
	resp, err := s.DB.Get(ctx, fullKey).Bytes()
	if err != nil {
		return nil, err
	}
	var schema fosite.DefaultClient
	if err = json.Unmarshal(resp, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

func (s FositeRedisStore) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	expTimeScore, err := s.DB.ZScore(ctx, prefixJTIBlocklist, jti).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get JTI score: %w", err)
	}
	currentTime := float64(time.Now().Unix())
	if expTimeScore > currentTime {
		return fosite.ErrJTIKnown
	}
	return nil
}

func (s FositeRedisStore) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	// TODO this could be pipelined
	currentTimeScore := float64(time.Now().Unix())
	s.DB.ZRemRangeByScore(ctx, prefixJTIBlocklist, "-inf", fmt.Sprintf("%f", currentTimeScore))
	rank, err := s.DB.ZRank(ctx, prefixJTIBlocklist, jti).Result()
	if err != nil {
		return errors.Wrap(err, "failed to get JTI rank")
	}
	if rank > -1 {
		// the entry exists, so the JTI is known
		return fosite.ErrJTIKnown
	}
	score := float64(exp.Unix())
	return s.DB.ZAdd(ctx, prefixJTIBlocklist, redis.Z{Score: score, Member: jti}).Err()
}

func (s FositeRedisStore) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {
	return s.setRequest(ctx, s.redisKey(prefixCode), code, req)
}

func (s FositeRedisStore) GetAuthorizeCodeSession(ctx context.Context, code string, sess fosite.Session) (fosite.Requester, error) {
	session, active, err := s.getRequest(ctx, s.redisKey(prefixCode), code, sess)
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	} else if !active {
		return nil, errors.Wrap(fosite.ErrInvalidatedAuthorizeCode, "")
	}

	return session, nil
}

func (s FositeRedisStore) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	return s.deactivateRequest(ctx, prefixCode, code)
}

func (s FositeRedisStore) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	sess, _, err := s.getRequest(ctx, s.redisKey(prefixPKCE), signature, session)
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return sess, nil
}

func (s FositeRedisStore) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.setRequest(ctx, s.redisKey(prefixPKCE), signature, requester)
}

func (s FositeRedisStore) DeletePKCERequestSession(ctx context.Context, signature string) error {
	return s.deleteRequest(ctx, s.redisKey(prefixPKCE), signature)
}

func (s FositeRedisStore) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	return s.redisCreateTokenSession(
		ctx,
		req,
		s.redisKey(prefixAccess),
		s.redisKey(prefixAccess, req.GetID()),
		signature,
	)
}

func (s FositeRedisStore) GetAccessTokenSession(ctx context.Context, signature string, sess fosite.Session) (fosite.Requester, error) {
	session, _, err := s.getRequest(ctx, s.redisKey(prefixAccess), signature, sess)
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return session, nil
}

func (s FositeRedisStore) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	return s.deleteRequest(ctx, s.redisKey(prefixAccess), signature)
}

func (s FositeRedisStore) CreateRefreshTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	return s.redisCreateTokenSession(
		ctx,
		req,
		s.redisKey(prefixRefresh),
		s.redisKey(prefixRefresh, req.GetID()),
		signature,
	)
}

func (s FositeRedisStore) GetRefreshTokenSession(ctx context.Context, signature string, sess fosite.Session) (fosite.Requester, error) {
	session, active, err := s.getRequest(ctx, s.redisKey(prefixRefresh), signature, sess)
	if err == redis.Nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	} else if !active {
		return nil, errors.Wrap(fosite.ErrInactiveToken, "")
	}
	return session, nil
}

func (s FositeRedisStore) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	return s.deleteRequest(ctx, s.redisKey(prefixRefresh), signature)
}

func (s FositeRedisStore) CreateImplicitAccessTokenSession(ctx context.Context, code string, req fosite.Requester) error {
	return s.CreateAccessTokenSession(ctx, code, req)
}

func (s FositeRedisStore) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, req fosite.Requester) error {
	if err := s.DB.Del(ctx, s.redisKey(prefixCode, authorizeCode)).Err(); err != nil {
		return err
	} else if err = s.CreateAccessTokenSession(ctx, accessSignature, req); err != nil {
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

func (s FositeRedisStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, req fosite.Requester) error {
	if err := s.DeleteRefreshTokenSession(ctx, originalRefreshSignature); err != nil {
		return err
	} else if err = s.CreateAccessTokenSession(ctx, accessSignature, req); err != nil {
		return err
	} else if err = s.CreateRefreshTokenSession(ctx, refreshSignature, req); err != nil {
		return err
	}
	return nil
}

func (s FositeRedisStore) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	return s.DB.Del(ctx, s.redisKey(prefixOIDC, authorizeCode)).Err()
}

func (s FositeRedisStore) RevokeRefreshToken(ctx context.Context, id string) error {
	refreshSet := s.redisKey(prefixRefresh, id)
	iter := s.DB.SScan(ctx, refreshSet, 0, "", 500).Iterator()
	sigs := make([]interface{}, 0)
	refreshKeys := make([]string, 0)
	for iter.Next(ctx) {
		sig := iter.Val()
		sigs = append(sigs, sig)
		refreshKeys = append(refreshKeys, s.redisKey(prefixRefresh, sig))
	}
	if err := iter.Err(); err != nil {
		return err
	}
	// delete each sig found in a loop. can't do the single DEL command because clustering
	// this could be optimized using MasterForKey to break the list into a single DEL command per shard, then doing
	// those concurrently
	for _, key := range refreshKeys {
		if err := s.DB.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return s.DB.SRem(ctx, refreshSet, sigs...).Err()
}

func (s FositeRedisStore) RevokeAccessToken(ctx context.Context, id string) error {
	accessSet := s.redisKey(prefixAccess, id)
	iter := s.DB.SScan(ctx, accessSet, 0, "", 500).Iterator()
	sigs := make([]interface{}, 0)
	refreshKeys := make([]string, 0)
	for iter.Next(ctx) {
		sig := iter.Val()
		sigs = append(sigs, sig)
		refreshKeys = append(refreshKeys, s.redisKey(prefixAccess, sig))
	}
	if err := iter.Err(); err != nil {
		return err
	}
	// delete each sig found in a loop. can't do the single DEL command because clustering
	// this could be optimized using MasterForKey to break the list into a single DEL command per shard, then doing
	// those concurrently
	for _, key := range refreshKeys {
		if err := s.DB.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return s.DB.SRem(ctx, accessSet, sigs...).Err()
}

func (s FositeRedisStore) GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (*jose.JSONWebKey, error) {
	// We actually don't want to use redis for this. Implementation from fosite MemoryStore:
	//if issuerKeys, ok := s.IssuerPublicKeys[issuer]; ok {
	//	if subKeys, ok := issuerKeys.KeysBySub[subject]; ok {
	//		if keyScopes, ok := subKeys.Keys[keyId]; ok {
	//			return keyScopes.Key, nil
	//		}
	//	}
	//}
	//
	//return nil, fosite.ErrNotFound
	panic("implement me")
}

func (s FositeRedisStore) GetPublicKeys(ctx context.Context, issuer string, subject string) (*jose.JSONWebKeySet, error) {
	// We actually don't want to use redis for this. Implementation from fosite MemoryStore:
	//if issuerKeys, ok := s.IssuerPublicKeys[issuer]; ok {
	//	if subKeys, ok := issuerKeys.KeysBySub[subject]; ok {
	//		if len(subKeys.Keys) == 0 {
	//			return nil, fosite.ErrNotFound
	//		}
	//
	//		keys := make([]jose.JSONWebKey, 0, len(subKeys.Keys))
	//		for _, keyScopes := range subKeys.Keys {
	//			keys = append(keys, *keyScopes.Key)
	//		}
	//
	//		return &jose.JSONWebKeySet{Keys: keys}, nil
	//	}
	//}
	//
	//return nil, fosite.ErrNotFound
	panic("implement me")
}

func (s FositeRedisStore) GetPublicKeyScopes(ctx context.Context, issuer string, subject string, keyId string) ([]string, error) {
	// same as GetPublicKey above, but return the Scopes field instead of Key
	panic("implement me")
}

func (s FositeRedisStore) IsJWTUsed(ctx context.Context, jti string) (bool, error) {
	err := s.ClientAssertionJWTValid(ctx, jti)
	if err != nil {
		return true, nil
	}
	return false, nil
}

func (s FositeRedisStore) MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) error {
	return s.SetClientAssertionJWT(ctx, jti, exp)
}

func (s FositeRedisStore) FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	// will not implement - this is only for the janitor command to clean up expired tokens. we should use redis TTLs for this
	return nil
}

func (s FositeRedisStore) FlushInactiveLoginConsentRequests(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	// will not implement - this is only for the janitor command to clean up expired tokens. we should use redis TTLs for this
	return nil
}

func (s FositeRedisStore) DeleteAccessTokens(ctx context.Context, clientID string) error {
	// this is supposed to delete all access tokens for a given client
	// no matter what, this is an expensive operation... we could do a large key scan...
	return nil
}

func (s FositeRedisStore) FlushInactiveRefreshTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	// will not implement - this is only for the janitor command to clean up expired tokens. we should use redis TTLs for this
	return nil
}

func (s FositeRedisStore) redisCreateTokenSession(ctx context.Context, req fosite.Requester, key, setKey, signature string) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = s.DB.Set(ctx, s.redisKey(key, signature), string(payload), 0).Err()
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = s.DB.SAdd(ctx, setKey, signature).Err()
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s FositeRedisStore) getRequest(ctx context.Context, prefix, key string, sess fosite.Session) (*fosite.Request, bool, error) {
	fullKey := s.redisKey(prefix, key)
	resp, err := s.DB.Get(ctx, fullKey).Bytes()
	if err != nil {
		return nil, false, err
	}
	var schema redisSchema
	if err = json.Unmarshal(resp, &schema); err != nil {
		return nil, false, err
	}
	return &fosite.Request{
		ID:                schema.ID,
		RequestedAt:       schema.RequestedAt,
		Client:            schema.Client,
		RequestedScope:    schema.RequestedScope,
		GrantedScope:      schema.GrantedScope,
		Form:              schema.Form,
		Session:           sess,
		RequestedAudience: schema.RequestedAudience,
		GrantedAudience:   schema.GrantedAudience,
		Lang:              schema.Lang,
	}, schema.Active, nil
}

func (s FositeRedisStore) setRequest(ctx context.Context, prefix, key string, requester fosite.Requester) error {
	payload, err := json.Marshal(requester)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fullKey := s.redisKey(prefix, key)
	if err = s.DB.Set(ctx, fullKey, string(payload), 0).Err(); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s FositeRedisStore) deleteRequest(ctx context.Context, prefix, signature string) error {
	sess, _, err := s.getRequest(ctx, prefix, signature, nil)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "")
	}
	err = s.DB.Del(ctx, s.redisKey(prefix, signature)).Err()
	if err != nil {
		return err
	}
	if sess != nil {
		err = s.DB.SRem(ctx, s.redisKey(prefix, sess.GetID()), signature).Err()
	}
	if err != nil {
		return err
	}
	return nil
}

func (s FositeRedisStore) deactivateRequest(ctx context.Context, prefix, key string) error {
	fullKey := s.redisKey(prefix, key)
	// WATCH/EXEC applies optimistic locking - the Set will fail if the key is modified while we're in the func
	return s.DB.Watch(ctx, func(tx *redis.Tx) error {
		resp, err := tx.Get(ctx, fullKey).Bytes()
		if err == redis.Nil {
			return fosite.ErrNotFound
		}
		if err != nil {
			return err
		}
		//var schema fosite.Request
		var schema redisSchema
		if err = json.Unmarshal(resp, &schema); err != nil {
			return err
		}
		schema.Active = false
		updatedSession, err := json.Marshal(schema)
		if err != nil {
			return err
		}
		err = tx.Set(ctx, fullKey, updatedSession, 0).Err()
		if err != nil {
			return err
		}
		return nil
	}, fullKey)
}

func (s FositeRedisStore) redisKey(fields ...string) string {
	return s.KeyPrefix + strings.Join(fields, ":")
}
