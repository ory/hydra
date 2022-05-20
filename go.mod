module github.com/ory/hydra

go 1.16

replace (
	github.com/bradleyjkemp/cupaloy/v2 => github.com/aeneasr/cupaloy/v2 v2.6.1-0.20210924214125-3dfdd01210a3
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/gobuffalo/packr => github.com/gobuffalo/packr v1.30.1
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/luna-duclos/instrumentedsql => github.com/ory/instrumentedsql v1.2.0
	github.com/luna-duclos/instrumentedsql/opentracing => github.com/ory/instrumentedsql/opentracing v0.0.0-20210903114257-c8963b546c5c
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.9
	github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1
)

require (
	github.com/ThalesIgnite/crypto11 v1.2.4
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-openapi/errors v0.20.1
	github.com/go-openapi/runtime v0.20.0
	github.com/go-openapi/strfmt v0.20.3
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.1
	github.com/go-swagger/go-swagger v0.26.1
	github.com/gobuffalo/pop/v6 v6.0.1
	github.com/gobuffalo/x v0.0.0-20181007152206-913e47c59ca7
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.0
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/luna-duclos/instrumentedsql/opentracing v0.0.0-20201103091713-40d03108b6f4
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/goveralls v0.0.11 // indirect
	github.com/miekg/pkcs11 v1.0.3
	github.com/mikefarah/yq/v4 v4.19.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.1
	github.com/olekukonko/tablewriter v0.0.1
	github.com/ory/analytics-go/v4 v4.0.3
	github.com/ory/fosite v0.42.2
	github.com/ory/go-acc v0.2.8
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.13
	github.com/ory/x v0.0.368
	github.com/pborman/uuid v1.2.1
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/cors v1.8.0
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.14.0
	github.com/toqueteos/webbrowser v1.2.0
	github.com/urfave/negroni v1.0.0
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/tools v0.1.10
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.33.0
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
)
