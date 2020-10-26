package clients

import (
	"bytes"
	"github.com/ory/hydra/internal/httpclient/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadClientFile(t *testing.T) {
	t.Run("reads client from json file", func(t *testing.T) {
		inputClient := map[string]interface{}{
			"client_id":   "test id",
			"client_name": "test name",
		}
		r := bytes.NewReader(
			[]byte(requireMarshaledJSON(t, inputClient)),
		)

		c, err := clientFromFile(r)
		require.NoError(t, err)

		assert.Equal(t, inputClient["client_id"], c.ClientID)
		assert.Equal(t, inputClient["client_name"], c.ClientName)
	})
}

func TestReadClientFlags(t *testing.T) {
	newTestCmd := func() (*cobra.Command, *bytes.Buffer) {
		stdErr := &bytes.Buffer{}
		cmd := &cobra.Command{}
		cmd.SetErr(stdErr)
		registerClientFlags(cmd.Flags())
		return cmd, stdErr
	}

	clientWithDefautls := func() *models.OAuth2Client {
		return &models.OAuth2Client{
			GrantTypes:              models.StringSlicePipeDelimiter{"authorization_code"},
			ResponseTypes:           models.StringSlicePipeDelimiter{"code"},
			AllowedCorsOrigins:      models.StringSlicePipeDelimiter{},
			PostLogoutRedirectUris:  models.StringSlicePipeDelimiter{},
			RedirectUris:            models.StringSlicePipeDelimiter{},
			Audience:                models.StringSlicePipeDelimiter{},
			SubjectType:             "public",
			TokenEndpointAuthMethod: "client_secret_basic",
		}
	}

	t.Run("reads client with defaults", func(t *testing.T) {
		cmd, stdErr := newTestCmd()

		expected := clientWithDefautls()

		actual, err := clientFromFlags(cmd)
		require.NoError(t, err)
		require.Equal(t, 0, stdErr.Len())

		assert.Equal(t, expected, actual)
	})

	t.Run("reads client with some flags", func(t *testing.T) {
		cmd, stdErr := newTestCmd()

		expected := clientWithDefautls()
		expected.ClientID = "test id 1392e9237401"
		expected.ClientName = "some name"
		expected.Audience = models.StringSlicePipeDelimiter{"aud1", "aud2"}

		require.NoError(t, cmd.Flags().Set(FlagClientID, expected.ClientID))
		require.NoError(t, cmd.Flags().Set(FlagClientName, expected.ClientName))
		require.NoError(t, cmd.Flags().Set(FlagAudience, strings.Join(expected.Audience, ",")))

		actual, err := clientFromFlags(cmd)
		require.NoError(t, err)
		require.Equal(t, 0, stdErr.Len())

		assert.Equal(t, expected, actual)
	})
}

func TestReadClientAllSources(t *testing.T) {
	newTestCmd := func() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
		stdIn, stdErr := &bytes.Buffer{}, &bytes.Buffer{}
		cmd := &cobra.Command{}
		cmd.SetIn(stdIn)
		cmd.SetErr(stdErr)
		registerClientFlags(cmd.Flags())
		return cmd, stdIn, stdErr
	}

	t.Run("reads from STD_IN when filename is -", func(t *testing.T) {
		cmd, stdIn, stdErr := newTestCmd()
		inputClient := map[string]interface{}{
			"client_id": "some id",
		}
		_, err := stdIn.WriteString(requireMarshaledJSON(t, inputClient))
		require.NoError(t, err)

		c, err := clientFromAllSources(cmd, "-")
		require.NoError(t, err)
		require.Equal(t, 0, stdErr.Len())

		assert.Equal(t, c.ClientID, inputClient["client_id"])
	})

	t.Run("reads from file", func(t *testing.T) {
		cmd, _, stdErr := newTestCmd()
		fn := filepath.Join(t.TempDir(), "client.json")
		inputClient := map[string]interface{}{
			"client_name": "can't think of a name...",
		}
		require.NoError(t, ioutil.WriteFile(fn, []byte(requireMarshaledJSON(t, inputClient)), 0600))

		c, err := clientFromAllSources(cmd, fn)
		require.NoError(t, err)
		require.Equal(t, 0, stdErr.Len())

		assert.Equal(t, inputClient["client_name"], c.ClientName)
	})

	t.Run("flags take precedence over files", func(t *testing.T) {
		cmd, stdIn, stdErr := newTestCmd()
		inputClient := map[string]interface{}{
			"client_id": "lahcvaliuwhei",
			"client_name": "liaiewnjch",
			"audience": []string{},
		}
		_, err := stdIn.WriteString(requireMarshaledJSON(t, inputClient))
		require.NoError(t, err)

		// overwrite audience here
		inputClient["audience"] = []string{"audx", "audy"}
		require.NoError(t, cmd.Flags().Set(FlagAudience, strings.Join(inputClient["audience"].([]string), ",")))

		c, err := clientFromAllSources(cmd, "-")
		require.NoError(t, err)
		require.Equal(t, 0, stdErr.Len())

		assert.Equal(t, inputClient["client_id"], c.ClientID)
		assert.Equal(t, inputClient["client_name"], c.ClientName)
		assert.Equal(t, models.StringSlicePipeDelimiter(inputClient["audience"].([]string)), c.Audience)
	})
}
