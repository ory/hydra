package jwk_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	. "github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/dancannon/gorethink.v2"

	"log"
	"os"
	"time"

	"github.com/ory-am/fosite/rand"
	"github.com/square/go-jose"
	"golang.org/x/net/context"
)

var managers = map[string]Manager{}

var testGenerator = &RS256Generator{}

var ts *httptest.Server

func init() {
	localWarden, httpClient := internal.NewFirewall(
		"tests",
		"alice",
		fosite.Arguments{
			"hydra.keys.create",
			"hydra.keys.get",
			"hydra.keys.delete",
			"hydra.keys.update",
		}, &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:keys:<faz|bar|foo><.*>"},
			Actions:   []string{"create", "get", "delete", "update"},
			Effect:    ladon.AllowAccess,
		},
	)

	r := httprouter.New()
	h := Handler{
		Manager: &MemoryManager{},
		W:       localWarden,
		H:       &herodot.JSON{},
	}
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	u, _ := url.Parse(ts.URL + "/keys")
	managers["memory"] = &MemoryManager{}
	managers["http"] = &HTTPManager{Client: httpClient, Endpoint: u}
}

var rethinkManager = new(RethinkManager)

func TestMain(m *testing.M) {
	var session *r.Session
	var err error

	c, err := dockertest.ConnectToRethinkDB(20, time.Second, func(url string) bool {
		if session, err = r.Connect(r.ConnectOpts{Address: url, Database: "hydra"}); err != nil {
			return false
		} else if _, err = r.DBCreate("hydra").RunWrite(session); err != nil {
			log.Printf("Database exists: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_keys").RunWrite(session); err != nil {
			log.Printf("Could not create table: %s", err)
			return false
		}

		key, err := rand.RandomBytes(32)
		if err != nil {
			log.Printf("Could not watch: %s", err)
			return false
		}
		rethinkManager = &RethinkManager{
			Keys:    map[string]jose.JsonWebKeySet{},
			Session: session,
			Table:   r.Table("hydra_keys"),
			Cipher: &AEAD{
				Key: key,
			},
		}
		rethinkManager.Watch(context.Background())
		time.Sleep(100 * time.Millisecond)
		return true
	})
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	managers["rethink"] = rethinkManager

	retCode := m.Run()
	c.KillRemove()
	os.Exit(retCode)
}

func BenchmarkRethinkGet(b *testing.B) {
	b.StopTimer()

	m := rethinkManager

	var err error
	ks, _ := testGenerator.Generate("")
	err = m.AddKeySet("newfoobar", ks)
	if err != nil {
		b.Fatalf("%s", err)
	}
	time.Sleep(time.Millisecond * 100)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.GetKey("newfoobar", "public")
	}
}

func TestColdStart(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	p1 := ks.Key("private")

	ks, _ = testGenerator.Generate("")
	p2 := ks.Key("private")

	pkg.AssertError(t, false, rethinkManager.AddKey("foo", First(p1)))
	pkg.AssertError(t, false, rethinkManager.AddKey("bar", First(p2)))

	time.Sleep(time.Second / 2)
	rethinkManager.Lock()
	rethinkManager.Keys = make(map[string]jose.JsonWebKeySet)
	rethinkManager.Unlock()
	pkg.AssertError(t, false, rethinkManager.ColdStart())

	c1, err := rethinkManager.GetKey("foo", "private")
	pkg.AssertError(t, false, err)
	c2, err := rethinkManager.GetKey("bar", "private")
	pkg.AssertError(t, false, err)

	assert.NotEqual(t, c1, c2)
	rethinkManager.Lock()
	rethinkManager.Keys = make(map[string]jose.JsonWebKeySet)
	rethinkManager.Unlock()
}

func TestManagerKey(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	priv := ks.Key("private")
	pub := ks.Key("public")

	for name, m := range managers {
		t.Logf("Running test %s", name)

		_, err := m.GetKey("faz", "baz")
		pkg.AssertError(t, true, err, name)

		err = m.AddKey("faz", First(priv))
		pkg.AssertError(t, false, err, name)

		time.Sleep(time.Millisecond * 100)

		got, err := m.GetKey("faz", "private")
		pkg.RequireError(t, false, err, name)
		assert.Equal(t, priv, got.Keys, "%s", name)

		err = m.AddKey("faz", First(pub))
		pkg.AssertError(t, false, err, name)

		time.Sleep(time.Millisecond * 100)

		got, err = m.GetKey("faz", "private")
		pkg.RequireError(t, false, err, name)
		assert.Equal(t, priv, got.Keys, "%s", name)

		got, err = m.GetKey("faz", "public")
		pkg.RequireError(t, false, err, name)
		assert.Equal(t, pub, got.Keys, "%s", name)

		err = m.DeleteKey("faz", "public")
		pkg.AssertError(t, false, err, name)

		time.Sleep(time.Millisecond * 100)

		ks, err = m.GetKey("faz", "public")
		pkg.AssertError(t, true, err, name)
	}

	err := managers["http"].AddKey("nonono", First(priv))
	pkg.AssertError(t, true, err, "%s")
}

func TestManagerKeySet(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	ks.Key("private")

	for name, m := range managers {
		_, err := m.GetKeySet("foo")
		pkg.AssertError(t, true, err, name)

		err = m.AddKeySet("bar", ks)
		pkg.AssertError(t, false, err, name)

		time.Sleep(time.Millisecond * 100)

		got, err := m.GetKeySet("bar")
		pkg.RequireError(t, false, err, name)
		assert.Equal(t, ks.Key("public"), got.Key("public"), name)
		assert.Equal(t, ks.Key("private"), got.Key("private"), name)

		err = m.DeleteKeySet("bar")
		pkg.AssertError(t, false, err, name)

		time.Sleep(time.Millisecond * 100)

		_, err = m.GetKeySet("bar")
		pkg.AssertError(t, true, err, name)
	}

	err := managers["http"].AddKeySet("nonono", ks)
	pkg.AssertError(t, true, err, "%s")
}
