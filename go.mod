module github.com/ory/hydra

go 1.16

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.7-0.20210414154423-1157a4212dcb

replace github.com/seatgeek/logrus-gelf-formatter => github.com/zepatrik/logrus-gelf-formatter v0.0.0-20210305135027-b8b3731dba10

replace github.com/dgrijalva/jwt-go => github.com/form3tech-oss/jwt-go v3.2.1+incompatible

replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2

replace github.com/gobuffalo/packr => github.com/gobuffalo/packr v1.30.1
replace github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1

require (
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/evanphx/json-patch v0.5.2
	github.com/go-bindata/go-bindata v3.1.1+incompatible
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/swag v0.19.13
	github.com/go-openapi/validate v0.20.1
	github.com/go-swagger/go-swagger v0.26.1
	github.com/gobuffalo/pop/v5 v5.3.4
	github.com/gobuffalo/x v0.0.0-20181007152206-913e47c59ca7
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.5.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.0
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/jackc/pgx/v4 v4.11.0
	github.com/jmoiron/sqlx v1.3.3
	github.com/julienschmidt/httprouter v1.3.0
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/luna-duclos/instrumentedsql/opentracing v0.0.0-20201103091713-40d03108b6f4
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.1
	github.com/olekukonko/tablewriter v0.0.1
	github.com/ory/analytics-go/v4 v4.0.1
	github.com/ory/fosite v0.40.2
	github.com/ory/go-acc v0.2.6
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.6
	github.com/ory/x v0.0.264
	github.com/pborman/uuid v1.2.1
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/rs/cors v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20190212093014-1007f53448d7
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.7.1
	github.com/toqueteos/webbrowser v1.2.0
	github.com/urfave/negroni v1.0.0
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	golang.org/x/tools v0.1.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.27.1
	gopkg.in/square/go-jose.v2 v2.5.1
)
