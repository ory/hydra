package oauth2

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

func (o *Handler) DefaultConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logrus.Warnln("It looks like no consent endpoint was set. All OAuth2 flows except client credentials will fail.")

	w.Write([]byte(`
<html>
<head>
	<title>Misconfigured consent endpoint</title>
</head>
<body>
<p>
	It looks like you forgot to set the consent endpoint url, which can be set using the <code>CONSENT_URL</code>
	environment variable.
</p>
<p>
	If you are an administrator, please read <a href="https://ory-am.gitbooks.io/hydra/content/oauth2.html">
	the guide</a> to understand what you need to do. If you are a user, please contact the administrator.
</p>
</body>
</html>
`))
}
