package oauth2

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RDBItems map[string]*RdbSchema

type FositeRehinkDBStore struct {
	Session *r.Session
	sync.RWMutex

	AuthorizeCodesTable r.Term
	IDSessionsTable     r.Term
	AccessTokensTable   r.Term
	RefreshTokensTable  r.Term
	ClientsTable        r.Term

	client.Manager

	AuthorizeCodes RDBItems
	IDSessions     RDBItems
	AccessTokens   RDBItems
	RefreshTokens  RDBItems
}

type RdbSchema struct {
	ID            string           `json:"id" gorethink:"id"`
	RequestID     string           `json:"requestId" gorethink:"requestId"`
	RequestedAt   time.Time        `json:"requestedAt" gorethink:"requestedAt"`
	Client        *client.Client   `json:"client" gorethink:"client"`
	Scopes        fosite.Arguments `json:"scopes" gorethink:"scopes"`
	GrantedScopes fosite.Arguments `json:"grantedScopes" gorethink:"grantedScopes"`
	Form          url.Values       `json:"form" gorethink:"form"`
	Session       json.RawMessage  `json:"session" gorethink:"session"`
}

func requestFromRDB(s *RdbSchema, proto fosite.Session) (*fosite.Request, error) {
	if proto != nil {
		if err := json.Unmarshal(s.Session, proto); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	d := new(fosite.Request)
	d.ID = s.RequestID
	d.RequestedAt = s.RequestedAt
	d.Client = s.Client
	d.Scopes = s.Scopes
	d.GrantedScopes = s.GrantedScopes
	d.Form = s.Form
	d.Session = proto
	return d, nil
}

func (m *FositeRehinkDBStore) ColdStart() error {
	if err := m.AccessTokens.coldStart(m.Session, &m.RWMutex, m.AccessTokensTable); err != nil {
		return err
	} else if err := m.AuthorizeCodes.coldStart(m.Session, &m.RWMutex, m.AuthorizeCodesTable); err != nil {
		return err
	} else if err := m.IDSessions.coldStart(m.Session, &m.RWMutex, m.IDSessionsTable); err != nil {
		return err
	} else if err := m.RefreshTokens.coldStart(m.Session, &m.RWMutex, m.RefreshTokensTable); err != nil {
		return err
	}
	return nil
}

func (s *FositeRehinkDBStore) publishInsert(table r.Term, id string, requester fosite.Requester) error {
	sess, err := json.Marshal(requester.GetSession())
	if err != nil {
		return errors.Wrap(err, "")
	}

	if _, err := table.Insert(&RdbSchema{
		ID:            id,
		RequestID:     requester.GetID(),
		RequestedAt:   requester.GetRequestedAt(),
		Client:        requester.GetClient().(*client.Client),
		Scopes:        requester.GetRequestedScopes(),
		GrantedScopes: requester.GetGrantedScopes(),
		Form:          requester.GetRequestForm(),
		Session:       sess,
	}).RunWrite(s.Session); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s *FositeRehinkDBStore) publishDelete(table r.Term, id string) error {
	if _, err := table.Get(id).Delete().RunWrite(s.Session); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (s *FositeRehinkDBStore) waitFor(i RDBItems, id string) error {
	c := make(chan bool)

	go func() {
		loopWait := time.Millisecond
		s.RLock()
		_, ok := i[id]
		s.RUnlock()
		for !ok {
			time.Sleep(loopWait)
			loopWait = loopWait * time.Duration(int64(2))
			if loopWait > time.Second {
				loopWait = time.Second
			}
			s.RLock()
			_, ok = i[id]
			s.RUnlock()
		}

		c <- true
	}()

	select {
	case <-c:
		return nil
	case <-time.After(time.Minute / 2):
		return errors.New("Timed out waiting for write confirmation")
	}
}

func (s *FositeRehinkDBStore) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	if err := s.publishInsert(s.IDSessionsTable, authorizeCode, requester); err != nil {
		return err
	}
	return s.waitFor(s.IDSessions, authorizeCode)
}

func (s *FositeRehinkDBStore) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	cl, ok := s.IDSessions[authorizeCode]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return requestFromRDB(cl, requester.GetSession())
}

func (s *FositeRehinkDBStore) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	return s.publishDelete(s.IDSessionsTable, authorizeCode)
}

func (s *FositeRehinkDBStore) CreateAuthorizeCodeSession(_ context.Context, code string, requester fosite.Requester) error {
	if err := s.publishInsert(s.AuthorizeCodesTable, code, requester); err != nil {
		return err
	}
	return s.waitFor(s.AuthorizeCodes, code)
}

func (s *FositeRehinkDBStore) GetAuthorizeCodeSession(_ context.Context, code string, sess fosite.Session) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.AuthorizeCodes[code]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	return requestFromRDB(rel, sess)
}

func (s *FositeRehinkDBStore) DeleteAuthorizeCodeSession(_ context.Context, code string) error {
	return s.publishDelete(s.AuthorizeCodesTable, code)
}

func (s *FositeRehinkDBStore) CreateAccessTokenSession(_ context.Context, signature string, requester fosite.Requester) error {
	if err := s.publishInsert(s.AccessTokensTable, signature, requester); err != nil {
		return err
	}
	return s.waitFor(s.AccessTokens, signature)
}

func (s *FositeRehinkDBStore) GetAccessTokenSession(_ context.Context, signature string, sess fosite.Session) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.AccessTokens[signature]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	return requestFromRDB(rel, sess)
}

func (s *FositeRehinkDBStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	return s.publishDelete(s.AccessTokensTable, signature)
}

func (s *FositeRehinkDBStore) CreateRefreshTokenSession(_ context.Context, signature string, requester fosite.Requester) error {
	if err := s.publishInsert(s.RefreshTokensTable, signature, requester); err != nil {
		return err
	}
	return s.waitFor(s.RefreshTokens, signature)
}

func (s *FositeRehinkDBStore) GetRefreshTokenSession(_ context.Context, signature string, sess fosite.Session) (fosite.Requester, error) {
	s.RLock()
	defer s.RUnlock()
	rel, ok := s.RefreshTokens[signature]
	if !ok {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	return requestFromRDB(rel, sess)
}

func (s *FositeRehinkDBStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	return s.publishDelete(s.RefreshTokensTable, signature)
}

func (s *FositeRehinkDBStore) CreateImplicitAccessTokenSession(ctx context.Context, code string, req fosite.Requester) error {
	return s.CreateAccessTokenSession(ctx, code, req)
}

func (s *FositeRehinkDBStore) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := s.DeleteAuthorizeCodeSession(ctx, authorizeCode); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	}

	if refreshSignature == "" {
		return nil
	}

	if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (s *FositeRehinkDBStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := s.DeleteRefreshTokenSession(ctx, originalRefreshSignature); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (m *FositeRehinkDBStore) Watch(ctx context.Context) {
	m.AccessTokens.watch(ctx, m.Session, &m.RWMutex, m.AccessTokensTable)
	m.AuthorizeCodes.watch(ctx, m.Session, &m.RWMutex, m.AuthorizeCodesTable)
	m.IDSessions.watch(ctx, m.Session, &m.RWMutex, m.IDSessionsTable)
	m.RefreshTokens.watch(ctx, m.Session, &m.RWMutex, m.RefreshTokensTable)
}

func (items RDBItems) coldStart(sess *r.Session, lock *sync.RWMutex, table r.Term) error {
	rows, err := table.Run(sess)
	if err != nil {
		return errors.Wrap(err, "")
	}

	var item RdbSchema
	lock.Lock()
	defer lock.Unlock()
	for rows.Next(&item) {
		var cp = item
		items[item.ID] = &cp
	}

	if rows.Err() != nil {
		return errors.Wrap(rows.Err(), "")
	}
	return nil
}

func (items RDBItems) watch(ctx context.Context, sess *r.Session, lock *sync.RWMutex, table r.Term) {
	go pkg.Retry(time.Second*15, time.Minute, func() error {
		lock.Lock()
		changes, err := table.Changes().Run(sess)
		lock.Unlock()
		if err != nil {
			return errors.Wrap(err, "")
		}
		defer changes.Close()

		var update = map[string]*RdbSchema{}
		for changes.Next(&update) {
			lock.Lock()
			logrus.Debugln("Received update from RethinkDB Cluster in OAuth2 manager.")
			newVal := update["new_val"]
			oldVal := update["old_val"]
			if newVal == nil && oldVal != nil {
				delete(items, oldVal.ID)
			} else if newVal != nil && oldVal != nil {
				delete(items, oldVal.ID)
				items[newVal.ID] = newVal
			} else {
				items[newVal.ID] = newVal
			}
			lock.Unlock()
		}

		if changes.Err() != nil {
			return errors.Wrap(changes.Err(), "")
		}

		return nil
	})
}

func (s *FositeRehinkDBStore) RevokeRefreshToken(ctx context.Context, id string) error {
	var found bool
	for sig, token := range s.RefreshTokens {
		if token.RequestID == id {
			if err := s.DeleteRefreshTokenSession(ctx, sig); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		return errors.New("Not found")
	}
	return nil
}

func (s *FositeRehinkDBStore) RevokeAccessToken(ctx context.Context, id string) error {
	var found bool
	for sig, token := range s.AccessTokens {
		if token.RequestID == id {
			if err := s.DeleteAccessTokenSession(ctx, sig); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		return errors.New("Not found")
	}
	return nil
}
