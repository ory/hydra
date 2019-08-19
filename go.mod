module github.com/ory/hydra

require (
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-openapi/analysis v0.19.0 // indirect
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/go-openapi/loads v0.19.0 // indirect
	github.com/go-openapi/runtime v0.19.0
	github.com/go-openapi/spec v0.19.0 // indirect
	github.com/go-openapi/strfmt v0.19.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-swagger/go-swagger v0.19.0
	github.com/go-swagger/scan-repo-boundary v0.0.0-20180623220736-973b3573c013 // indirect
	github.com/gobuffalo/packd v0.0.0-20190315124812-a385830c7fc0 // indirect
	github.com/gobuffalo/packr v1.24.0
	github.com/gobwas/glob v0.2.3
	github.com/golang/gddo v0.0.0-20190312205958-5a2505f3dbf0 // indirect
	github.com/golang/mock v1.3.1
	github.com/gorilla/handlers v1.4.0 // indirect
	github.com/gorilla/sessions v1.1.4-0.20181208214519-12bd4761fc66
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/imdario/mergo v0.0.0-20171009183408-7fe0c75c13ab
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kr/pty v1.1.3 // indirect
	github.com/lib/pq v1.0.0
	github.com/luna-duclos/instrumentedsql v0.0.0-20190316074304-ecad98b20aec // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/meatballhat/negroni-logrus v0.0.0-20170801195057-31067281800f
	github.com/mendsley/gojwk v0.0.0-20141217222730-4d5ec6e58103
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/ory/fosite v0.29.7
	github.com/ory/go-acc v0.0.0-20181118080137-ddc355013f90
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.6.0
	github.com/ory/x v0.0.68
	github.com/pborman/uuid v1.2.0
	github.com/phayes/freeport v0.0.0-20171002181615-b8543db493a5
	github.com/pkg/errors v0.8.1
	github.com/pkg/profile v1.3.0 // indirect
	github.com/prometheus/client_golang v0.9.3
	github.com/rs/cors v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20190212093014-1007f53448d7
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518
	github.com/stretchr/testify v1.3.0
	github.com/toqueteos/webbrowser v0.0.0-20150720201625-21fc9f95c834
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/tools v0.0.0-20190711191110-9a621aea19f8
	gopkg.in/square/go-jose.v2 v2.3.0
)

// Fix for https://github.com/golang/lint/issues/436
replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
