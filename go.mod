module github.com/ory/hydra

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-bindata/go-bindata v3.1.1+incompatible
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.5
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.3
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-swagger/go-swagger v0.21.1-0.20200107003254-1c98855b472d
	github.com/gobuffalo/packd v0.0.0-20190315124812-a385830c7fc0 // indirect
	github.com/gobuffalo/packr v1.24.0
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/sessions v1.1.4-0.20181208214519-12bd4761fc66
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.2.0
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mendsley/gojwk v0.0.0-20141217222730-4d5ec6e58103
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/opentracing/opentracing-go v1.1.1-0.20190913142402-a7454ce5950e
	github.com/ory/fosite v0.30.2
	github.com/ory/go-acc v0.0.0-20181118080137-ddc355013f90
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.6.2
	github.com/ory/sdk/swagutil v0.0.0-20200131170418-ead0c2285f93
	github.com/ory/viper v1.5.6
	github.com/ory/x v0.0.93
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/phayes/freeport v0.0.0-20171002181615-b8543db493a5
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.3.0 // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/rs/cors v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20190212093014-1007f53448d7
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518
	github.com/stretchr/testify v1.4.0
	github.com/toqueteos/webbrowser v1.2.0
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/urfave/negroni v1.0.0
	go.opentelemetry.io/otel v0.2.1
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/tools v0.0.0-20200203023011-6f24f261dadb
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

// Fix for https://github.com/golang/lint/issues/436
replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1

go 1.13
