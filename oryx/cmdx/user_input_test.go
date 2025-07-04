// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAskForConfirmation(t *testing.T) {
	t.Run("case=prints question", func(t *testing.T) {
		testQuestion := "test-question"
		stdin, stdout := new(bytes.Buffer), new(bytes.Buffer)

		_, err := stdin.Write([]byte("y\n"))
		require.NoError(t, err)

		AskForConfirmation(testQuestion, stdin, stdout)

		prompt, err := io.ReadAll(stdout)
		require.NoError(t, err)
		assert.Contains(t, string(prompt), testQuestion)
	})

	t.Run("case=accept", func(t *testing.T) {
		for _, input := range []string{
			"y\n",
			"yes\n",
		} {
			stdin := new(bytes.Buffer)

			_, err := stdin.Write([]byte(input))
			require.NoError(t, err)

			confirmed := AskForConfirmation("", stdin, new(bytes.Buffer))

			assert.True(t, confirmed)
		}
	})

	t.Run("case=reject", func(t *testing.T) {
		for _, input := range []string{
			"n\n",
			"no\n",
		} {
			stdin := new(bytes.Buffer)

			_, err := stdin.Write([]byte(input))
			require.NoError(t, err)

			confirmed := AskForConfirmation("", stdin, new(bytes.Buffer))

			assert.False(t, confirmed)
		}
	})

	t.Run("case=reprompt on random input", func(t *testing.T) {
		testQuestion := "question"

		for _, input := range []string{
			"foo\ny\n",
			"bar\nn\n",
		} {
			stdin, stdout := new(bytes.Buffer), new(bytes.Buffer)

			_, err := stdin.Write([]byte(input))
			require.NoError(t, err)

			AskForConfirmation(testQuestion, stdin, stdout)

			output, err := io.ReadAll(stdout)
			require.NoError(t, err)
			assert.Equal(t, 2, bytes.Count(output, []byte(testQuestion)))
		}
	})
}
