package oauth2

import (
	"fmt"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var consentManagers = map[string]ConsentRequestManager{
	"memory": NewConsentRequestMemoryManager(),
}

func TestConsentRequestManagerReadWrite(t *testing.T) {
	req := &ConsentRequest{ID: uuid.New()}
	for k, m := range consentManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetConsentRequest("1234")
			assert.Error(t, err)

			require.NoError(t, m.PersistConsentRequest(req))

			got, err := m.GetConsentRequest(req.ID)
			require.NoError(t, err)

			assert.EqualValues(t, req, got)
		})
	}
}

func TestConsentRequestManagerUpdate(t *testing.T) {
	req := &ConsentRequest{ID: uuid.New()}
	for k, m := range consentManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			require.NoError(t, m.PersistConsentRequest(req))

			got, err := m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.False(t, got.IsConsentGranted())

			require.NoError(t, m.AcceptConsentRequest(req.ID, new(AcceptConsentRequestPayload)))
			got, err = m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.True(t, got.IsConsentGranted())

			require.NoError(t, m.RejectConsentRequest(req.ID))
			got, err = m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.False(t, got.IsConsentGranted())
		})
	}
}

func TestHttpRequestClient(t *testing.T) {
	req := &ConsentRequest{ID: uuid.New()}
	for k, m := range consentManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			require.NoError(t, m.PersistConsentRequest(req))

			got, err := m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.False(t, got.IsConsentGranted())

			require.NoError(t, m.AcceptConsentRequest(req.ID, new(AcceptConsentRequestPayload)))
			got, err = m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.True(t, got.IsConsentGranted())

			require.NoError(t, m.RejectConsentRequest(req.ID))
			got, err = m.GetConsentRequest(req.ID)
			require.NoError(t, err)
			assert.False(t, got.IsConsentGranted())
		})
	}
}
