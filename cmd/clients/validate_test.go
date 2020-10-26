package clients

import (
	"bytes"
	"github.com/ory/x/cmdx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateCmd(t *testing.T) {
	t.Run("validates from STD_IN", func(t *testing.T) {
		cmd := newValidateCmd()
		input := map[string]interface{}{
			"client_id":   "pasdjfha",
			"client_name": "öavichna",
			"audience":    []string{"kljhsdcas", "ljca"},
		}
		stdIn := bytes.NewBufferString(requireMarshaledJSON(t, input))

		stdOut, stdErr, err := cmdx.Exec(t, cmd, stdIn, "-")
		require.NoError(t, err)
		require.Len(t, stdErr, 0)

		assert.Contains(t, stdOut, "valid")
	})
}
