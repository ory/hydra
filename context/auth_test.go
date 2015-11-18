package context

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ory-am/dockertest"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/policy"
	ladon "github.com/ory-am/ladon/policy/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var db *sql.DB
var ladonStore *ladon.Store

func TestMain(m *testing.M) {
	var err error
	var c dockertest.ContainerID
	c, db, err = dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	ladonStore = ladon.New(db)
	if err := ladonStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	os.Exit(m.Run())
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
			//			"exp": "2099-10-31T15:03:52.4620974+01:00",
			//			"iat": "2014-10-31T13:03:52.4620974+01:00",
			//			"nbf": "2014-10-31T13:03:52.4620974+01:00",
			//			"sub": "132"
			//		}
			"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDk5LTEwLTMxVDE1OjAzOjUyLjQ2MjA5NzQrMDE6MDAiLCJpYXQiOiIyMDE0LTEwLTMxVDEzOjAzOjUyLjQ2MjA5NzQrMDE6MDAiLCJuYmYiOiIyMDE0LTEwLTMxVDEzOjAzOjUyLjQ2MjA5NzQrMDE6MDAiLCJzdWIiOiIxMzIifQ.qnZr-msiG5GkVTDTyY3g26c5Edho36_E9CaANyCBVOrXWRfRPDMf7E2vrdZubO5tXlfKRgM_1avFQVWZhqrdrGBO8DiBa5OGX9IdAZaclqQFjg7vRSyIFllSs4zP4QREG4YL0qwiYGKS4SBcCS2LNfbaJfrKP_zUReXRAlWNdeFAw6zsGzlAtHQO_O0HnJCEB_wEBIkMIxdI2f-1yyTZJInyvY_wrFDkCkTfkmmW8EHzO2R44FXmaudxDCG1YAeN6WssAwgzBjR8WaQ2M_8VUYWN9TCDc3Fx58XWRTtWL_coDI9R6WtqaPkyr2_qn1Un3y3yLCGdVglRYnhJL1YCXA",
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
