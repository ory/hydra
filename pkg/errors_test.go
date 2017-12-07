package pkg

import (
	"testing"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"bytes"
	"strings"
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
