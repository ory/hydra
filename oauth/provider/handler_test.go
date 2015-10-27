package provider_test

import (
	"bytes"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/ory-am/dockertest"
	acpg "github.com/ory-am/hydra/account/postgres"
	"github.com/ory-am/hydra/hash"
	"github.com/ory-am/hydra/jwt"
	. "github.com/ory-am/hydra/oauth/provider"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	ppg "github.com/ory-am/ladon/policy/postgres"
	opg "github.com/ory-am/osin-storage/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

var db *sql.DB
var accountStore *acpg.Store
var policyStore *ppg.Store
var handler *Handler
var gd = new(guard.Guard)
var osinStore *opg.Storage

func TestMain(m *testing.M) {
	var err error
	var c dockertest.ContainerID
	c, db, err = dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	accountStore = acpg.New(&hash.BCrypt{10}, db)
	policyStore = ppg.New(db)
	osinStore = opg.New(db)
	if err := policyStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	if err := accountStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	if err := osinStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	j := jwt.New([]byte(jwt.TestCertificates[0][1]), []byte(jwt.TestCertificates[1][1]))
	handler = NewHandler(osinStore, j, accountStore, policyStore, gd)

	if _, err := osinStore.CreateClient("1", "secret", "/callback"); err != nil {
		log.Fatalf("Could create client: %s", err)
	}
	if _, err := accountStore.Create("2", "2@bar.com", "secret", "{}"); err != nil {
		log.Fatalf("Could create account: %s", err)
	}
	if _, err := policyStore.Create("3", "", policy.AllowAccess, []string{}, []string{"authorize"}, []string{"/oauth2/authorize"}); err != nil {
		log.Fatalf("Could create client: %s", err)
	}
	os.Exit(m.Run())
}

func TestAuthorize(t *testing.T) {
	oauthConfigs := []*oauth2.Config{
		&oauth2.Config{
			ClientID:     "1",
			ClientSecret: "secret",
			Scopes:       []string{},
			RedirectURL:  "/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "/oauth2/auth",
				TokenURL: "/oauth2/token",
			},
		},
		&oauth2.Config{
			ClientID:     "1",
			ClientSecret: "wrongsecret",
			Scopes:       []string{},
			RedirectURL:  "/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "/oauth2/auth",
				TokenURL: "/oauth2/token",
			},
		},
		&oauth2.Config{
			ClientID:     "notexistent",
			ClientSecret: "random",
			Scopes:       []string{},
			RedirectURL:  "/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "/oauth2/auth",
				TokenURL: "/oauth2/token",
			},
		},
	}

	type userData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type test struct {
		code     int
		state    string
		config   *oauth2.Config
		userData *userData
		pass     bool
	}

	cases := []*test{
		&test{
			state:    "foobar",
			config:   oauthConfigs[0],
			code:     http.StatusFound,
			userData: &userData{"2@bar.com", "secret"},
			pass:     true,
		},
		&test{
			state:    "foobar",
			config:   oauthConfigs[0],
			code:     http.StatusUnauthorized,
			userData: &userData{"nonexistent@bar.com", "secret"},
			pass:     false,
		},
		&test{
			state:    "foobar",
			config:   oauthConfigs[0],
			code:     http.StatusUnauthorized,
			userData: &userData{"2@bar.com", "wrong secret"},
			pass:     false,
		},
		&test{
			state:  "foobar",
			config: oauthConfigs[1],
			// Ok because oauth2/auth does not check client secret, only oauth2/token does.
			code:     http.StatusFound,
			userData: &userData{"2@bar.com", "secret"},
			pass:     false,
		},
		&test{
			state:    "foobar",
			config:   oauthConfigs[2],
			code:     http.StatusUnauthorized,
			userData: &userData{"2@bar.com", "secret"},
			pass:     false,
		},
	}

	for k, c := range cases {
		loc := ""
		func() {
			router := mux.NewRouter()
			handler.SetRoutes(router)

			authURL := c.config.AuthCodeURL(c.state)
			log.Printf("Acquired auth code url: %s", authURL)

			post := url.Values{}
			post.Set("username", c.userData.Username)
			post.Add("password", c.userData.Password)
			req, _ := http.NewRequest("POST", authURL, bytes.NewBufferString(post.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(t, c.code, res.Code, `Case %d, %s: %s`, k, res.Body.Bytes())

			log.Printf("Result was: %s %s", res.Body.String(), res.Header().Get("Location"))
			loc = res.Header().Get("Location")
		}()

		if loc == "" {
			continue
		}

		func() {
			router := mux.NewRouter()
			handler.SetRoutes(router)
			ts := httptest.NewServer(router)
			defer ts.Close()
			u, err := url.Parse(loc)
			require.Nil(t, err)
			log.Printf("Exchanging token: %s", ts.URL+"/oauth2/auth")
			c.config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}
			tok, err := c.config.Exchange(oauth2.NoContext, u.Query().Get("code"))
			if !c.pass {
				assert.NotNil(t, err, "Case %d", k)
				return
			}

			assert.Nil(t, err, "Case %d: %v", k)
			assert.NotNil(t, tok)
		}()
	}
}
