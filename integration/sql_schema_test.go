package integration

import (
	"testing"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	lsql "github.com/ory/ladon/manager/sql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSQLSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
		return
	}

	var testGenerator = &jwk.RS256Generator{}
	ks, _ := testGenerator.Generate("")
	p1 := ks.Key("private")
	r := fosite.NewRequest()
	r.ID = "foo"
	db := ConnectToPostgres()

	cm := &client.SQLManager{DB: db, Hasher: &fosite.BCrypt{}}
	gm := group.SQLManager{DB: db}
	jm := jwk.SQLManager{DB: db, Cipher: &jwk.AEAD{Key: []byte("11111111111111111111111111111111")}}
	om := oauth2.FositeSQLStore{Manager: cm, DB: db, L: logrus.New()}
	crm := oauth2.NewConsentRequestSQLManager(db)
	pm := lsql.NewSQLManager(db, nil)

	_, err := pm.CreateSchemas("", "hydra_ladon_migration")
	require.NoError(t, err)
	_, err = cm.CreateSchemas()
	require.NoError(t, err)
	_, err = gm.CreateSchemas()
	require.NoError(t, err)
	_, err = jm.CreateSchemas()
	require.NoError(t, err)
	_, err = om.CreateSchemas()
	require.NoError(t, err)
	_, err = crm.CreateSchemas()
	require.NoError(t, err)

	require.NoError(t, jm.AddKey("foo", jwk.First(p1)))
	require.NoError(t, pm.Create(&ladon.DefaultPolicy{ID: "foo"}))
	require.NoError(t, cm.CreateClient(&client.Client{ID: "foo"}))
	require.NoError(t, crm.PersistConsentRequest(&oauth2.ConsentRequest{ID: "foo"}))
	require.NoError(t, om.CreateAccessTokenSession(nil, "asdfasdf", r))
	require.NoError(t, gm.CreateGroup(&group.Group{
		ID:      "asdfas",
		Members: []string{"asdf"},
	}))
}
