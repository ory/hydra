package flowcache

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/x/sqlcon"
)

type contextKey int

const (
	contextKeyCache contextKey = iota
)

var lruCache *lru.Cache[string, *cacheEntry]

func init() {
	var err error
	lruCache, err = lru.New[string, *cacheEntry](10000)
	if err != nil {
		panic(err)
	}
}

type (
	Entry interface {
		Client() *client.Client
		Flow() *flow.Flow
		LoginSession() *consent.LoginSession
		Set(value any) (stored bool)

		PersistFlow(ctx context.Context, provider persistence.Persister) error
	}

	cacheEntry struct {
		state        string
		client       *client.Client
		flow         *flow.Flow
		loginSession *consent.LoginSession
		persistOnce  sync.Once
		persistErr   error
	}

	nopEntry struct{}
)

var nopEntryVar = new(nopEntry)

func (n *nopEntry) Client() *client.Client                                       { return nil }
func (n *nopEntry) Flow() *flow.Flow                                             { return nil }
func (n *nopEntry) LoginSession() *consent.LoginSession                          { return nil }
func (n *nopEntry) Set(_ any) bool                                               { return false }
func (n *nopEntry) PersistFlow(_ context.Context, _ persistence.Persister) error { return nil }

func (c *cacheEntry) Client() *client.Client              { return c.client }
func (c *cacheEntry) Flow() *flow.Flow                    { return c.flow }
func (c *cacheEntry) LoginSession() *consent.LoginSession { return c.loginSession }
func (c *cacheEntry) Set(value any) bool {
	switch v := value.(type) {
	case *client.Client:
		c.client = v
	case *flow.Flow:
		c.flow = v
	case *consent.LoginSession:
		c.loginSession = v
	default:
		panic(fmt.Sprintf("unexpected type %T", value))
	}
	return true
}
func (c *cacheEntry) PersistFlow(ctx context.Context, provider persistence.Persister) error {
	c.persistOnce.Do(func() {
		if c.flow != nil {
			c.persistErr = sqlcon.HandleError(provider.Connection(ctx).Create(c.flow))
		}
		lruCache.Remove(c.state)
	})
	return c.persistErr
}

func getEntry(state string) *cacheEntry {
	entry, ok := lruCache.Get(state)
	if !ok {
		entry = new(cacheEntry)
		entry.state = state
		lruCache.Add(state, entry)
	}

	return entry
}

func FromRequest(ctx context.Context, r *http.Request) (context.Context, Entry) {
	state := r.URL.Query().Get("state")
	entry := getEntry(state)

	return context.WithValue(ctx, contextKeyCache, entry), entry
}

func FromChallenge(ctx context.Context, challenge string) (context.Context, Entry) {
	// FIXME: POC, This is very inefficient
	for _, key := range lruCache.Keys() {
		entry := getEntry(key)
		if entry.Flow() == nil {
			continue
		}

		if entry.Flow().ID == challenge {
			return context.WithValue(ctx, contextKeyCache, entry), entry
		}
	}

	return ctx, nil
}

func FromConsentChallenge(ctx context.Context, consentChallenge string) (context.Context, Entry) {
	// FIXME: POC, This is very inefficient
	for _, key := range lruCache.Keys() {
		entry := getEntry(key)
		if entry.Flow() == nil {
			continue
		}

		if entry.Flow().ConsentChallengeID.String() == consentChallenge {
			return context.WithValue(ctx, contextKeyCache, entry), entry
		}
	}

	return ctx, nil
}

func FromContext(ctx context.Context) (entry Entry) {
	val := ctx.Value(contextKeyCache)
	if val == nil {
		return nopEntryVar
	}

	return val.(*cacheEntry)
}
