package context

import (
	"database/sql"
	"fmt"
	"github.com/ory-am/dockertest"
	hydra "github.com/ory-am/hydra/account/postgres"
	"github.com/ory-am/hydra/hash"
	"github.com/ory-am/hydra/jwt"
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
var hydraStore *hydra.Store
var ladonStore *ladon.Store

func TestMain(m *testing.M) {
	var err error
	var c dockertest.ContainerID
	c, db, err = dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	hydraStore = hydra.New(&hash.BCrypt{10}, db)
	ladonStore = ladon.New(db)
	if err := ladonStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	if err := hydraStore.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}
	os.Exit(m.Run())
}

type test struct {
	id              string
	privateKey      []byte
	publicKey       []byte
	authorization   string
	isAuthenticated bool
}

var cases = []test{
	test{
		"1",
		[]byte(jwt.TestCertificates[0][1]),
		[]byte(jwt.TestCertificates[1][1]),
		// {"foo": "bar"}
		"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		false,
	},
	test{
		"2",
		[]byte(jwt.TestCertificates[0][1]),
		[]byte(jwt.TestCertificates[1][1]),
		// {"subject": "nonexistent"}
		"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWJqZWN0Ijoibm9uZXhpc3RlbnQifQ.jDUnvVMQHrhuIRUr8qAJ0g-ZKArdiJ21LAPDktmV56KFknX712Yxdder78YjEjxvGOvgtxLpCiay0cV5pvcWLuFW65Ys1P1SwdmdebtWfiGQwBy2Ggm3MrHjD_-r5JNAxFZjFZfZ1Fk-JlSZ97r8S7gYfDSAkxhpDmDy5Bm8e5_xsGDNp8dByuXop7QEtJb_igaa0APWa2ZOp3oTgxjD4CP6ZX6N5fGjtwjJWx5wHt7JaKXq8CRG8elm7LnNezYyJxeHECVctQGVv3HUjJxKf0l7wZXbG87BrG2M7otT8Py2sJP8X4wYL0DEsbErkEieV4D-KEBqpkvfXOrDGMFNRQ",
		false,
	},
	test{
		"3",
		[]byte(jwt.TestCertificates[0][1]),
		[]byte(jwt.TestCertificates[1][1]),
		// not a valid token
		"Bearer eyJ0eXAaOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWJqZWN0IjoiMTMyIn0.WDC51GK5wIQPd00MqLjFy3AU3gNsvsCpWk5e8RObVxBqcYAdv-UwMfEUAFE6Y50C5pQ1t8_LHfxJYNfcW3fj_x5FXckdbqvpXHi-psxuDwk_rancpjZQegcutqYRH37_lnJ8lIq65ZgxnyYnQKGOMl3w7etK1gOvqEcP_eHn8HG0jeVk0SDZm82x0JXSk0lrVEEjWmWYtXEsLz0E4clNPUW37K9eyjYFKnyVCIPfmGwTlkDLjANsyu0P6kFiV28V1_XedtJXDI3MmG2SxSHogDhZJLb298JBwod0d6wTyygI9mUbX-C0PklTJTxIhSs7Pc6unNlWnbyL8Z4FJrdSEw",
		false,
	},
	test{
		"4",
		[]byte(jwt.TestCertificates[0][1]),
		[]byte(jwt.TestCertificates[1][1]),
		// {"subject": "132"}
		"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWJqZWN0IjoiMTIzIn0.cKlYJu0LPQHtIuHC1rdLEVIaSG-2OHJWC9zkDnt1Uh1W1aQz6VI2SJKriUkEq4vHQIG8V3tbuBGqzKK376twJmzN7uvxRP8e5-3Lewfp7mx6A1_vmnBXYuCaUazeDoZBWkQ3SDqmvg2kWnFJyp9VJtZm2fMojH76wnlyQCAjNaIOPOxNH0AqjdbMR6n9FdXOxJDZy9fYYR7qA-HJNSOYfr1fakbSNrvLeh3vjN95-854JuGSpmrDT78sz_opqF67ZUSpgONKElpTT5yLk3MVwxt9zTYCR4IjEkh23C6iqm25DDsATPZmLOt8XTiIJ-5tNa-J1vDPO8jD2__cC-Gu4A",
		true,
	},
	test{
		"4",
		[]byte(jwt.TestCertificates[0][1]),
		[]byte(jwt.TestCertificates[1][1]),
		"",
		false,
	},
}

func TestNewContextFromAuthorization(t *testing.T) {
	_, err := hydraStore.Create("123", "foo@bar", "secret", "{}")
	assert.Nil(t, err)
	for _, c := range cases {
		newContextFromAuthorization(t, c)
	}
}

func newContextFromAuthorization(t *testing.T, c test) {
	message := "ok"
	ctx := context.Background()

	j := jwt.New(c.privateKey, c.publicKey)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = NewContextFromAuthorization(ctx, r, j, hydraStore, ladonStore)
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
