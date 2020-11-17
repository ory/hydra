module github.com/ory/hydra

go 1.15

require (
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-bindata/go-bindata v3.1.1+incompatible
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/runtime v0.19.11
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.6
	github.com/go-swagger/go-swagger v0.22.1-0.20200306221957-4aad3a5f78b8
	github.com/gobuffalo/httptest v1.0.2
	github.com/gobuffalo/packr v1.24.0 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/gobuffalo/pop/v5 v5.3.1
	github.com/gobuffalo/x v0.0.0-20181007152206-913e47c59ca7
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.4.3
	github.com/google/uuid v1.1.1
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.0
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/jackc/pgx/v4 v4.9.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/luna-duclos/instrumentedsql/opentracing v0.0.0-20200611091901-487c5ec83473
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/ory/analytics-go/v4 v4.0.1
	github.com/ory/cli v0.0.28
	github.com/ory/fosite v0.36.0
	github.com/ory/go-acc v0.2.5
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.1
	github.com/ory/viper v1.7.5
	github.com/ory/x v0.0.163
	github.com/pborman/uuid v1.2.0
	github.com/phayes/freeport v0.0.0-20171002181615-b8543db493a5
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.3.0 // indirect
	github.com/prometheus/client_golang v1.4.0
	github.com/rs/cors v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20190212093014-1007f53448d7
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/toqueteos/webbrowser v1.2.0
	github.com/urfave/negroni v1.0.0
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/tools v0.0.0-20201019175715-b894a3290fff
	gopkg.in/DataDog/dd-trace-go.v1 v1.27.1
	gopkg.in/square/go-jose.v2 v2.5.1
)

replace github.com/gobuffalo/pop/v5 => github.com/gobuffalo/pop/v5 v5.3.2-0.20201029132236-f36afb546df1
