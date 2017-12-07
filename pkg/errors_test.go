package pkg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type errStackTracer struct{}

func (s *errStackTracer) StackTrace() errors.StackTrace {
	return errors.StackTrace{}
}

func (s *errStackTracer) Error() string {
	return "foo"
}

func TestLogError(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := logrus.New()
	l.Level = logrus.DebugLevel
	l.Out = buf
	LogError(errors.New("asdf"), l)

	t.Logf("%s", string(buf.Bytes()))

	assert.True(t, strings.Contains(string(buf.Bytes()), "Stack trace"))

	LogError(errors.Wrap(new(errStackTracer), ""), l)
}

func TestLogErrorDoesNotPanic(t *testing.T) {
	LogError(errors.New("asdf"), nil)
}
