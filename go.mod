module github.com/ory/hydra

go 1.16

replace (
	github.com/bradleyjkemp/cupaloy/v2 => github.com/aeneasr/cupaloy/v2 v2.6.1-0.20210924214125-3dfdd01210a3
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/gobuffalo/packr => github.com/gobuffalo/packr v1.30.1
	github.com/gobuffalo/pop/v6 => github.com/gobuffalo/pop/v6 v6.0.2-alpha-ci.0.20220421231416-e6ba76ba3be3
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.13
	github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1
)

replace github.com/ory/fosite => github.com/ory/fosite v0.42.3-0.20220617175535-a88d4431f12d

replace github.com/gobuffalo/pop/v6 => github.com/gobuffalo/pop/v6 v6.0.4-0.20220524160009-195240e4a669

require (
	github.com/ThalesIgnite/crypto11 v1.2.4
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/bradleyjkemp/cupaloy/v2 v2.6.0
	github.com/bxcodec/faker/v3 v3.7.0
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/fatih/structs v1.1.0
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-openapi/errors v0.20.1
	github.com/go-openapi/runtime v0.20.0
	github.com/go-openapi/strfmt v0.20.3
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.1
	github.com/go-swagger/go-swagger v0.26.1
	github.com/gobuffalo/pop/v6 v6.0.4-0.20220524160009-195240e4a669
	github.com/gobuffalo/x v0.0.0-20181007152206-913e47c59ca7
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.2 // indirect
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/hashicorp/go-retryablehttp v0.7.1
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/instana/testify v1.6.2-0.20200721153833-94b1851f4d65
	github.com/jackc/pgx/v4 v4.16.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/goveralls v0.0.11 // indirect
	github.com/miekg/pkcs11 v1.0.3
	github.com/mikefarah/yq/v4 v4.19.1
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.1
	github.com/olekukonko/tablewriter v0.0.1
	github.com/ory/analytics-go/v4 v4.0.3
	github.com/ory/fosite v0.42.3-0.20220513181618-5f156bd07d5d
	github.com/ory/go-acc v0.2.8
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.13
	github.com/ory/x v0.0.418
	github.com/pborman/uuid v1.2.1
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/cors v1.8.2
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.1
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/tidwall/gjson v1.14.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/sjson v1.2.4
	github.com/toqueteos/webbrowser v1.2.0
	github.com/urfave/negroni v1.0.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.32.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.7.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.7.0 // indirect
	go.opentelemetry.io/contrib/samplers/jaegerremote v0.2.0 // indirect
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.7.0 // indirect
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/crypto v0.0.0-20220517005047-85d78b3ac167
	golang.org/x/net v0.0.0-20220524220425-1d687d428aca // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/tools v0.1.10
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	google.golang.org/genproto v0.0.0-20220525015930-6ca3db687a9d // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/tools v0.1.10
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.38.0
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
)
