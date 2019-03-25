module github.com/ory/hydra

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-errors/errors v1.0.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.3.0 // indirect
	github.com/google/uuid v1.1.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.1.3
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/imdario/mergo v0.0.0-20171009183408-7fe0c75c13ab
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.0.0
	github.com/meatballhat/negroni-logrus v0.0.0-20170801195057-31067281800f
	github.com/mendsley/gojwk v0.0.0-20141217222730-4d5ec6e58103
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.0
	github.com/opentracing/opentracing-go v1.0.2
	github.com/ory/dockertest v3.3.4+incompatible
	github.com/ory/fosite v0.29.0
	github.com/ory/go-convenience v0.1.0
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.5.1
	github.com/ory/hive-cloud/hive v0.0.0-20190312103236-eedac51faab3
	github.com/ory/sqlcon v0.0.7
	github.com/ory/x v0.0.37
	github.com/pborman/uuid v1.2.0
	github.com/phayes/freeport v0.0.0-20171002181615-b8543db493a5
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/prometheus/common v0.2.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190225181712-6ed1f7e10411 // indirect
	github.com/rs/cors v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20190212093014-1007f53448d7
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	github.com/stretchr/testify v1.3.0
	github.com/toqueteos/webbrowser v0.0.0-20150720201625-21fc9f95c834
	github.com/uber/jaeger-client-go v2.15.0+incompatible
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20190228161510-8dd112bcdc25
	golang.org/x/net v0.0.0-20190227022144-312bce6e941f // indirect
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421
	golang.org/x/sys v0.0.0-20190305064518-30e92a19ae4a // indirect
	google.golang.org/genproto v0.0.0-20190226184841-fc2db5cae922 // indirect
	google.golang.org/grpc v1.19.0 // indirect
	gopkg.in/resty.v1 v1.9.1
	gopkg.in/square/go-jose.v2 v2.2.2
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
)

replace github.com/ory/x => ../x
