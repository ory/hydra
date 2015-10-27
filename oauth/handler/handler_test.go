package handler_test

import (
	"database/sql"
	"encoding/json"
	"github.com/RangelReale/osin"
	"github.com/gorilla/mux"
	"github.com/ory-am/dockertest"
	acpg "github.com/ory-am/hydra/account/postgres"
	"github.com/ory-am/hydra/hash"
	"github.com/ory-am/hydra/jwt"
	. "github.com/ory-am/hydra/oauth/handler"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	ppg "github.com/ory-am/ladon/policy/postgres"
	opg "github.com/ory-am/osin-storage/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
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

	if err := osinStore.CreateClient(&osin.DefaultClient{"1", "secret", "/callback", ""}); err != nil {
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

var configs = map[string]*oauth2.Config{
	"working": &oauth2.Config{
		ClientID:     "1",
		ClientSecret: "secret",
		Scopes:       []string{},
		RedirectURL:  "/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "/oauth2/auth",
			TokenURL: "/oauth2/token",
		},
	},
	"voidSecret": &oauth2.Config{
		ClientID:     "1",
		ClientSecret: "wrongsecret",
		Scopes:       []string{},
		RedirectURL:  "/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "/oauth2/auth",
			TokenURL: "/oauth2/token",
		},
	},
	"voidID": &oauth2.Config{
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

type userAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var logins = map[string]*userAuth{
	"working":      &userAuth{"2@bar.com", "secret"},
	"voidEmail":    &userAuth{"1@bar.com", "secret"},
	"voidPassword": &userAuth{"1@bar.com", "public"},
}

func TestPasswordGrantType(t *testing.T) {
	type test struct {
		config *oauth2.Config
		user   *userAuth
		pass   bool
	}
	router := mux.NewRouter()
	handler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, c := range []*test{
		&test{configs["working"], logins["working"], true},
		&test{configs["working"], logins["voidEmail"], false},
		&test{configs["working"], logins["voidPassword"], false},
		&test{configs["working"], logins["working"], true},
		&test{configs["voidSecret"], logins["working"], false},
		&test{configs["voidID"], logins["working"], false},
		&test{configs["working"], logins["working"], true},
	} {
		config := *c.config
		config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}
		_, err := config.PasswordCredentialsToken(oauth2.NoContext, c.user.Username, c.user.Password)
		if c.pass {
			assert.Nil(t, err, "Case %d", k)
		} else {
			assert.NotNil(t, err, "Case %d", k)
		}
	}
}

func TestClientGrantType(t *testing.T) {
	type test struct {
		config *oauth2.Config
		pass   bool
	}
	router := mux.NewRouter()
	handler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, c := range []*test{
		&test{configs["working"], true},
		&test{configs["voidSecret"], false},
		&test{configs["voidID"], false},
		&test{configs["working"], true},
	} {
		conf := clientcredentials.Config{
			ClientID:     c.config.ClientID,
			ClientSecret: c.config.ClientSecret,
			TokenURL:     ts.URL + c.config.Endpoint.TokenURL,
			Scopes:       c.config.Scopes,
		}

		_, err := conf.Token(oauth2.NoContext)
		if c.pass {
			assert.Nil(t, err, "Case %d\n%v", k, conf)
		} else {
			assert.NotNil(t, err, "Case %d\n%v", k, conf)
		}
	}
}

func TestIntrospect(t *testing.T) {
	router := mux.NewRouter()
	handler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	config := *configs["working"]
	user := logins["working"]
	clientConfig := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     ts.URL + config.Endpoint.TokenURL,
		Scopes:       config.Scopes,
	}
	config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}

	access, err := clientConfig.Token(oauth2.NoContext)
	require.Nil(t, err)
	verify, err := config.PasswordCredentialsToken(oauth2.NoContext, user.Username, user.Password)
	require.Nil(t, err)

	client := &http.Client{}
	form := url.Values{}
	form.Add("token", access.AccessToken)

	req, err := http.NewRequest("POST", ts.URL+"/oauth2/introspect", strings.NewReader(form.Encode()))
	require.Nil(t, err)

	req.Header.Add("Authorization", "Bearer "+verify.AccessToken)
	res, err := client.Do(req)
	require.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode, "%v", body)

	var result map[string]interface{}
	require.Nil(t, json.Unmarshal(body, &result), "%v\n%v", body)
	assert.True(t, result["active"].(bool))
	t.Logf("Got token: %s", body)
}

//
//
//func skipTestAuthorize(t *testing.T) {
//
//	type test struct {
//		code     int
//		state    string
//		config   *oauth2.Config
//		userData *userAuth
//		pass     bool
//	}
//
//	cases := []*test{
//		&test{
//			state:    "foobar",
//			config:   oauthConfigs[0],
//			code:     http.StatusFound,
//			userData: &userAuth{"2@bar.com", "secret"},
//			pass:     true,
//		},
//		&test{
//			state:    "foobar",
//			config:   oauthConfigs[0],
//			code:     http.StatusUnauthorized,
//			userData: &userAuth{"nonexistent@bar.com", "secret"},
//			pass:     false,
//		},
//		&test{
//			state:    "foobar",
//			config:   oauthConfigs[0],
//			code:     http.StatusUnauthorized,
//			userData: &userAuth{"2@bar.com", "wrong secret"},
//			pass:     false,
//		},
//		&test{
//			state:  "foobar",
//			config: oauthConfigs[1],
//			// Ok because oauth2/auth does not check client secret, only oauth2/token does.
//			code:     http.StatusFound,
//			userData: &userAuth{"2@bar.com", "secret"},
//			pass:     false,
//		},
//		&test{
//			state:    "foobar",
//			config:   oauthConfigs[2],
//			code:     http.StatusUnauthorized,
//			userData: &userAuth{"2@bar.com", "secret"},
//			pass:     false,
//		},
//	}
//
//	for k, c := range cases {
//		loc := ""
//
//		// FIXME This test case is actually what we don't want and it should be removed
//		func() {
//			router := mux.NewRouter()
//			handler.SetRoutes(router)
//
//			authURL := c.config.AuthCodeURL(c.state)
//			log.Printf("Acquired auth code url: %s", authURL)
//
//			// FIXME The oauth2/auth endpoint is not for authentication. This should not work
//			post := url.Values{}
//			post.Set("username", c.userData.Username)
//			post.Add("password", c.userData.Password)
//			req, _ := http.NewRequest("POST", authURL, bytes.NewBufferString(post.Encode()))
//			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//			res := httptest.NewRecorder()
//			router.ServeHTTP(res, req)
//			assert.Equal(t, c.code, res.Code, `Case %d, %s: %s`, k, res.Body.Bytes())
//
//			log.Printf("Result was: %s %s", res.Body.String(), res.Header().Get("Location"))
//			loc = res.Header().Get("Location")
//		}()
//
//		if loc == "" {
//			continue
//		}
//
//		func() {
//			router := mux.NewRouter()
//			handler.SetRoutes(router)
//			ts := httptest.NewServer(router)
//			defer ts.Close()
//			u, err := url.Parse(loc)
//			require.Nil(t, err)
//			log.Printf("Exchanging token: %s", ts.URL + "/oauth2/auth")
//			c.config.Endpoint = oauth2.Endpoint{AuthURL: ts.URL + "/oauth2/auth", TokenURL: ts.URL + "/oauth2/token"}
//			tok, err := c.config.Exchange(oauth2.NoContext, u.Query().Get("code"))
//			if !c.pass {
//				assert.NotNil(t, err, "Case %d", k)
//				return
//			}
//
//			assert.Nil(t, err, "Case %d: %v", k)
//			assert.NotNil(t, tok)
//		}()
//	}
//}
