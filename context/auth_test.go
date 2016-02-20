package context

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/policy"
	ladon "github.com/ory-am/ladon/policy/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gopkg.in/ory-am/dockertest.v2"
)

var db *sql.DB
var ladonStore *ladon.Store

func TestMain(m *testing.M) {
	var err error
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	ladonStore = ladon.New(db)
	if err := ladonStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	db.Close()
	c.KillRemove()
}

func TestNewContextFromAuthorization(t *testing.T) {
	for _, c := range []struct {
		id              string
		privateKey      []byte
		publicKey       []byte
		authorization   string
		isAuthenticated bool
	}{
		{
			"1",
			[]byte(hjwt.TestCertificates[0][1]),
			[]byte(hjwt.TestCertificates[1][1]),
			// {"foo": "bar"}
			"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			false,
		},
		{
			"2",
			[]byte(hjwt.TestCertificates[0][1]),
			[]byte(hjwt.TestCertificates[1][1]),
			// {"subject": "nonexistent"}
			"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWJqZWN0Ijoibm9uZXhpc3RlbnQifQ.jDUnvVMQHrhuIRUr8qAJ0g-ZKArdiJ21LAPDktmV56KFknX712Yxdder78YjEjxvGOvgtxLpCiay0cV5pvcWLuFW65Ys1P1SwdmdebtWfiGQwBy2Ggm3MrHjD_-r5JNAxFZjFZfZ1Fk-JlSZ97r8S7gYfDSAkxhpDmDy5Bm8e5_xsGDNp8dByuXop7QEtJb_igaa0APWa2ZOp3oTgxjD4CP6ZX6N5fGjtwjJWx5wHt7JaKXq8CRG8elm7LnNezYyJxeHECVctQGVv3HUjJxKf0l7wZXbG87BrG2M7otT8Py2sJP8X4wYL0DEsbErkEieV4D-KEBqpkvfXOrDGMFNRQ",
			false,
		},
		{
			"3",
			[]byte(hjwt.TestCertificates[0][1]),
			[]byte(hjwt.TestCertificates[1][1]),
			// not a valid token
			"Bearer eyJ0eXAaOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWJqZWN0IjoiMTMyIn0.WDC51GK5wIQPd00MqLjFy3AU3gNsvsCpWk5e8RObVxBqcYAdv-UwMfEUAFE6Y50C5pQ1t8_LHfxJYNfcW3fj_x5FXckdbqvpXHi-psxuDwk_rancpjZQegcutqYRH37_lnJ8lIq65ZgxnyYnQKGOMl3w7etK1gOvqEcP_eHn8HG0jeVk0SDZm82x0JXSk0lrVEEjWmWYtXEsLz0E4clNPUW37K9eyjYFKnyVCIPfmGwTlkDLjANsyu0P6kFiV28V1_XedtJXDI3MmG2SxSHogDhZJLb298JBwod0d6wTyygI9mUbX-C0PklTJTxIhSs7Pc6unNlWnbyL8Z4FJrdSEw",
			false,
		},
		{
			"4",
			[]byte(hjwt.TestCertificates[0][1]),
			[]byte(hjwt.TestCertificates[1][1]),
			//		{
			//			"exp": 1924975619,
			//			"iat": 1924975619,
			//			"nbf": 0
			//			"aud": "tests",
			//			"subject": "foo"
			//		}
			"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtYXgiLCJleHAiOjE5MjQ5NzU2MTksIm5iZiI6MCwiaWF0IjoxOTI0OTc1NjE5LCJhdWQiOiJ0ZXN0cyIsInN1YmplY3QiOiJmb28ifQ.lvjLGnLO3mZSS63fomK-KH2mhLXjjg9b13opiN7jY4MrXE_DaR0Lum8a_RcqqSTXbpHxYSIPV9Ji7zM_X1bvBtsPpBE1PR3_PrdD5_uIDQ-UWPVzozxhOvuZzU7qHx3TFQClZ6tYIXYioTszz9zQHiE4hj1x6Z_shWPfczELGyD0HnEC3o_w7IFfYO_L0YDN_vkuqr6yS5kaPIsoCF_iHuhTzoBAEIpUENlxSpCPuxR9aMaJ-BQDInHoPc1h-VvkgOdR_iENQdOUePObw17ywdGkRk6C5kRHSxjca-ULGcDn36NZ54SEPolcGbjs3vVA1g0jQARKIcTVw6Uu7x0s6Q",
			true,
		},
		{
			"5",
			[]byte(hjwt.TestCertificates[0][1]),
			[]byte(hjwt.TestCertificates[1][1]),
			"",
			false,
		},
	} {
		message := "ok"
		ctx := context.Background()

		j := hjwt.New(c.privateKey, c.publicKey)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx = NewContextFromAuthorization(ctx, r, j, ladonStore)
			assert.Equal(t, c.isAuthenticated, IsAuthenticatedFromContext(ctx), "Case %s", c.id)
			fmt.Fprintln(w, message)
		}))
		defer ts.Close()

		client := &http.Client{}
		req, err := http.NewRequest("GET", ts.URL, nil)
		require.Nil(t, err)
		req.Header.Set("Authorization", c.authorization)
		res, err := client.Do(req)
		require.Nil(t, err)

		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		require.Nil(t, err)
		assert.Equal(t, message+"\n", string(result))
	}
}

func TestGetters(t *testing.T) {
	assert.False(t, IsAuthenticatedFromContext(context.Background()))
	_, err := PoliciesFromContext(context.Background())
	assert.NotNil(t, err)
	_, err = SubjectFromContext(context.Background())
	assert.NotNil(t, err)
	_, err = TokenFromContext(context.Background())
	assert.NotNil(t, err)

	ctx := context.Background()
	claims := hjwt.ClaimsCarrier{"sub": "peter"}
	token := &jwt.Token{Valid: true}
	policies := []policy.Policy{}
	ctx = NewContextFromAuthValues(ctx, claims, token, policies)

	assert.True(t, IsAuthenticatedFromContext(ctx))
	policiesContext, err := PoliciesFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, policies, policiesContext)

	subjectContext, err := SubjectFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, claims.GetSubject(), subjectContext)

	tokenContext, err := TokenFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, token, tokenContext)
}
