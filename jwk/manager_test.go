package jwk_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/compose"
	"github.com/ory-am/hydra/herodot"
	. "github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/dancannon/gorethink.v2"

	"log"
	"os"
	"time"

	"crypto/rand"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

var managers = map[string]Manager{}

var testGenerator = &RS256Generator{}

var ts *httptest.Server
var httpManager *HTTPManager

func init() {
	localWarden, httpClient := compose.NewFirewall(
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
			Resources: []string{"rn:hydra:keys:<faz|bar|foo|anonymous><.*>"},
			Actions:   []string{"create", "get", "delete", "update"},
			Effect:    ladon.AllowAccess,
		}, &ladon.DefaultPolicy{
			ID:        "2",
			Subjects:  []string{"alice", ""},
			Resources: []string{"rn:hydra:keys:anonymous<.*>"},
			Actions:   []string{"get"},
			Effect:    ladon.AllowAccess,
		},
	)

	router := httprouter.New()
	h := Handler{
		Manager: &MemoryManager{},
		W:       localWarden,
		H:       &herodot.JSON{},
	}
	h.SetRoutes(router)
	ts := httptest.NewServer(router)
	u, _ := url.Parse(ts.URL + "/keys")
	managers["memory"] = &MemoryManager{}
	httpManager = &HTTPManager{Client: httpClient, Endpoint: u}
	managers["http"] = httpManager
}

var rethinkManager = new(RethinkManager)

func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.Wrap(err, "")
	}
	return bytes, nil
}

var encryptionKey, _ = randomBytes(32)

var containers = []dockertest.ContainerID{}

func TestMain(m *testing.M) {
	defer func() {
		for _, c := range containers {
			c.KillRemove()
		}
	}()

	connectToMySQL()
	connectToRethinkDB()
	connectToPG()

	os.Exit(m.Run())
}

func connectToPG() {
	var db *sqlx.DB
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sqlx.Open("postgres", url)
		if err != nil {
			log.Printf("Got error in postgres connector: %s", err)
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	containers = append(containers, c)
	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}

	if err = s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	managers["postgres"] = s
	containers = append(containers, c)
}

func connectToRethinkDB() {
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

		rethinkManager = &RethinkManager{
			Keys:    map[string]jose.JsonWebKeySet{},
			Session: session,
			Table:   r.Table("hydra_keys"),
			Cipher: &AEAD{
				Key: encryptionKey,
			},
		}
		rethinkManager.Watch(context.Background())
		time.Sleep(100 * time.Millisecond)
		return true
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	containers = append(containers, c)
	managers["rethink"] = rethinkManager
}

func connectToMySQL() {
	var db *sqlx.DB
	c, err := dockertest.ConnectToMySQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sqlx.Open("mysql", url)
		if err != nil {
			log.Printf("Got error in mysql connector: %s", err)
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	containers = append(containers, c)
	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}

	if err = s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	managers["mysql"] = s
	containers = append(containers, c)
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

func TestHTTPManagerPublicKeyGet(t *testing.T) {
	anonymous := &HTTPManager{Endpoint: httpManager.Endpoint, Client: http.DefaultClient}
	ks, _ := testGenerator.Generate("")
	priv := ks.Key("private")

	name := "http"
	m := httpManager

	_, err := m.GetKey("anonymous", "baz")
	pkg.AssertError(t, true, err, name)

	err = m.AddKey("anonymous", First(priv))
	pkg.AssertError(t, false, err, name)

	time.Sleep(time.Millisecond * 100)

	got, err := anonymous.GetKey("anonymous", "private")
	pkg.RequireError(t, false, err, name)
	assert.Equal(t, priv, got.Keys, "%s", name)
}

func TestManagerKey(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	priv := ks.Key("private")
	pub := ks.Key("public")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), func(t *testing.T) {
			_, err := m.GetKey("faz", "baz")
			assert.NotNil(t, err)

			err = m.AddKey("faz", First(priv))
			assert.Nil(t, err)

			time.Sleep(time.Millisecond * 100)

			got, err := m.GetKey("faz", "private")
			assert.Nil(t, err)
			assert.Equal(t, priv, got.Keys, "%s", name)

			err = m.AddKey("faz", First(pub))
			assert.Nil(t, err)

			time.Sleep(time.Millisecond * 100)

			got, err = m.GetKey("faz", "private")
			assert.Nil(t, err)
			assert.Equal(t, priv, got.Keys, "%s", name)

			got, err = m.GetKey("faz", "public")
			assert.Nil(t, err)
			assert.Equal(t, pub, got.Keys, "%s", name)

			err = m.DeleteKey("faz", "public")
			assert.Nil(t, err)

			time.Sleep(time.Millisecond * 100)

			ks, err = m.GetKey("faz", "public")
			assert.NotNil(t, err)
		})
	}

	err := managers["http"].AddKey("nonono", First(priv))
	pkg.AssertError(t, true, err, "%s")
}

func TestManagerKeySet(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	ks.Key("private")

	for name, m := range managers {
		t.Run(fmt.Sprintf("case=%s", name), func(t *testing.T) {
			_, err := m.GetKeySet("foo")
			pkg.AssertError(t, true, err, name)

			err = m.AddKeySet("bar", ks)
			assert.Nil(t, err)

			time.Sleep(time.Millisecond * 100)

			got, err := m.GetKeySet("bar")
			assert.Nil(t, err)
			assert.Equal(t, ks.Key("public"), got.Key("public"), name)
			assert.Equal(t, ks.Key("private"), got.Key("private"), name)

			err = m.DeleteKeySet("bar")
			assert.Nil(t, err)

			time.Sleep(time.Millisecond * 100)

			_, err = m.GetKeySet("bar")
			assert.NotNil(t, err)
		})
	}

	err := managers["http"].AddKeySet("nonono", ks)
	pkg.AssertError(t, true, err, "%s")
}
