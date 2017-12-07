package pkg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogError(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := logrus.New()
	l.Level = logrus.DebugLevel
	l.Out = buf
	LogError(errors.New("asdf"), l)

	t.Logf("%s", string(buf.Bytes()))

	assert.True(t, strings.Contains(string(buf.Bytes()), "Stack trace"))
}

func TestLogErrorDoesNotPanic(t *testing.T) {
	LogError(errors.New("asdf"), nil)
}
