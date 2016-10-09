package internal

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
	ImplicitTable       r.Term
	RefreshTokensTable  r.Term
	ClientsTable        r.Term

	client.Manager

	AuthorizeCodes RDBItems
	IDSessions     RDBItems
	AccessTokens   RDBItems
	Implicit       RDBItems
	RefreshTokens  RDBItems
}

type RdbSchema struct {
	ID            string           `json:"id" gorethink:"id"`
	RequestedAt   time.Time        `json:"requestedAt" gorethink:"requestedAt"`
	Client        *client.Client   `json:"client" gorethink:"client"`
	Scopes        fosite.Arguments `json:"scopes" gorethink:"scopes"`
	GrantedScopes fosite.Arguments `json:"grantedScopes" gorethink:"grantedScopes"`
	Form          url.Values       `json:"form" gorethink:"form"`
	Session       json.RawMessage  `json:"session" gorethink:"session"`
}

func requestFromRDB(s *RdbSchema, proto interface{}) (*fosite.Request, error) {
	if proto != nil {
		if err := json.Unmarshal(s.Session, proto); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	d := new(fosite.Request)
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
	} else if err := m.Implicit.coldStart(m.Session, &m.RWMutex, m.ImplicitTable); err != nil {
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

func waitFor(i RDBItems, id string) error {
	c := make(chan bool)

	go func() {
		loopWait := time.Millisecond
		_, ok := i[id]
		for !ok {
			time.Sleep(loopWait)
			loopWait = loopWait * time.Duration(int64(2))
			if loopWait > time.Second {
				loopWait = time.Second
			}
			_, ok = i[id]
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
	return waitFor(s.IDSessions, authorizeCode)
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
	return waitFor(s.AuthorizeCodes, code)
}

func (s *FositeRehinkDBStore) GetAuthorizeCodeSession(_ context.Context, code string, sess interface{}) (fosite.Requester, error) {
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
	return waitFor(s.AccessTokens, signature)
}

func (s *FositeRehinkDBStore) GetAccessTokenSession(_ context.Context, signature string, sess interface{}) (fosite.Requester, error) {
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
	return waitFor(s.RefreshTokens, signature)
}

func (s *FositeRehinkDBStore) GetRefreshTokenSession(_ context.Context, signature string, sess interface{}) (fosite.Requester, error) {
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

func (s *FositeRehinkDBStore) CreateImplicitAccessTokenSession(_ context.Context, code string, req fosite.Requester) error {
	if err := s.publishInsert(s.ImplicitTable, code, req); err != nil {
		return err
	}
	return waitFor(s.Implicit, code)
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
	m.Implicit.watch(ctx, m.Session, &m.RWMutex, m.ImplicitTable)
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
		changes, err := table.Changes().Run(sess)
		if err != nil {
			return errors.Wrap(err, "")
		}
		defer changes.Close()

		var update = map[string]*RdbSchema{}
		for changes.Next(&update) {
			logrus.Debugln("Received update from RethinkDB Cluster in OAuth2 manager.")
			newVal := update["new_val"]
			oldVal := update["old_val"]
			lock.Lock()
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
