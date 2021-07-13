package x

import (
	"net/http"

	"github.com/ory/x/logrusx"
)

func LogAudit(r *http.Request, message interface{}, logger *logrusx.Logger) {
	if logger == nil {
		logger = logrusx.NewAudit("", "")
	}

	logger = logger.WithRequest(r)

	if err, ok := message.(error); ok {
		logger.WithError(err).Infoln("access denied")
		return
	}

	logger.Infoln("access allowed")
}
