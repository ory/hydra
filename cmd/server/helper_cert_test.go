package server_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/x/tlsx"

	"github.com/ory/hydra/cmd/server"
	"github.com/ory/hydra/jwk"
)

func TestGetOrCreateTLSCertificate(t *testing.T) {
	keys, err := (&jwk.RS256Generator{KeyLength: 1024}). // this is ok because we're just using it for tests
								Generate(uuid.New().String(), "sig")
	require.NoError(t, err)

	private := keys.Keys[0]
	cert, err := tlsx.CreateSelfSignedCertificate(private.Key)
	require.NoError(t, err)
	server.AttachCertificate(&private, cert)

	var actual jose.JSONWebKeySet
	var b bytes.Buffer
	require.NoError(t, json.NewEncoder(&b).Encode(keys))
	require.NoError(t, json.NewDecoder(&b).Decode(&actual))
}
