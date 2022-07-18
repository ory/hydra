module github.com/ory/hydra

go 1.18

replace (
	github.com/bradleyjkemp/cupaloy/v2 => github.com/aeneasr/cupaloy/v2 v2.6.1-0.20210924214125-3dfdd01210a3
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/gobuffalo/packr => github.com/gobuffalo/packr v1.30.1
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.13
	github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1
)

replace github.com/ory/hydra-client-go => ./internal/httpclient

require (
	github.com/ThalesIgnite/crypto11 v1.2.4
	github.com/bradleyjkemp/cupaloy/v2 v2.6.0
	github.com/bxcodec/faker/v3 v3.7.0
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/fatih/structs v1.1.0
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-swagger/go-swagger v0.26.1
	github.com/gobuffalo/pop/v6 v6.0.4-0.20220707121509-feb0eed33d86
	github.com/gobuffalo/x v0.0.0-20181007152206-913e47c59ca7
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/hashicorp/go-retryablehttp v0.7.1
	github.com/instana/testify v1.6.2-0.20200721153833-94b1851f4d65
	github.com/jackc/pgx/v4 v4.16.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/miekg/pkcs11 v1.0.3
	github.com/mikefarah/yq/v4 v4.19.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oleiade/reflections v1.0.1
	github.com/ory/analytics-go/v4 v4.0.3
	github.com/ory/fosite v0.42.3-0.20220623093045-89faad818c4c
	github.com/ory/go-acc v0.2.8
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.13
	github.com/ory/hydra-client-go v1.11.8
	github.com/ory/jsonschema/v3 v3.0.7
	github.com/ory/x v0.0.443
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/cors v1.8.2
	github.com/sawadashota/encrypta v0.0.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	github.com/tidwall/gjson v1.14.0
	github.com/tidwall/sjson v1.2.4
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80
	github.com/toqueteos/webbrowser v1.2.0
	github.com/urfave/negroni v1.0.0
	go.opentelemetry.io/otel v1.7.0
	go.step.sm/crypto v0.16.2
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	golang.org/x/tools v0.1.10
	gopkg.in/DataDog/dd-trace-go.v1 v1.39.0
	gopkg.in/square/go-jose.v2 v2.6.0
)

require (
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.0.0-20211129110424-6491aa3bf583 // indirect
	github.com/DataDog/datadog-go v4.8.2+incompatible // indirect
	github.com/DataDog/datadog-go/v5 v5.0.2 // indirect
	github.com/DataDog/sketches-go v1.2.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/a8m/envsubst v1.3.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/avast/retry-go/v4 v4.0.5 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cockroachdb/cockroach-go/v2 v2.2.10 // indirect
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/dave/jennifer v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/docker/cli v20.10.14+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.9+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/ecordell/optgen v0.0.6 // indirect
	github.com/elastic/go-licenser v0.3.1 // indirect
	github.com/elastic/go-sysinfo v1.1.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/elliotchance/orderedmap v1.4.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.19.16 // indirect
	github.com/go-openapi/errors v0.20.1 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/loads v0.20.1 // indirect
	github.com/go-openapi/runtime v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-openapi/strfmt v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-openapi/validate v0.20.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gobuffalo/envy v1.10.1 // indirect
	github.com/gobuffalo/fizz v1.14.0 // indirect
	github.com/gobuffalo/flect v0.2.5 // indirect
	github.com/gobuffalo/github_flavored_markdown v1.1.1 // indirect
	github.com/gobuffalo/helpers v0.6.4 // indirect
	github.com/gobuffalo/nulls v0.4.1 // indirect
	github.com/gobuffalo/plush/v4 v4.1.12 // indirect
	github.com/gobuffalo/tags/v3 v3.1.2 // indirect
	github.com/gobuffalo/validate/v3 v3.3.1 // indirect
	github.com/goccy/go-yaml v1.9.5 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/inhies/go-bytesize v0.0.0-20210819104631-275770b98743 // indirect
	github.com/instana/go-sensor v1.41.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/jandelgado/gcov2lcov v1.0.5 // indirect
	github.com/jcchavezs/porto v0.1.0 // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/knadh/koanf v1.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/looplab/fsm v0.1.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/mattn/goveralls v0.0.11 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/microcosm-cc/bluemonday v1.0.16 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/nyaruka/phonenumbers v1.0.73 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.1.1 // indirect
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5 // indirect
	github.com/openzipkin/zipkin-go v0.4.0 // indirect
	github.com/ory/dockertest/v3 v3.9.0 // indirect
	github.com/ory/go-convenience v0.1.0 // indirect
	github.com/ory/viper v1.7.5 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/profile v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.4 // indirect
	github.com/seatgeek/logrus-gelf-formatter v0.0.0-20210414080842-5b05eb8ff761 // indirect
	github.com/segmentio/backo-go v0.0.0-20200129164019-23eae7c10bd3 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.8.1 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/thales-e-security/pool v0.0.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/timtadh/data-structures v0.5.3 // indirect
	github.com/timtadh/lexmachine v0.2.2 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	go.elastic.co/apm v1.15.0 // indirect
	go.elastic.co/apm/module/apmhttp v1.15.0 // indirect
	go.elastic.co/apm/module/apmot v1.15.0 // indirect
	go.elastic.co/fastjson v1.1.0 // indirect
	go.mongodb.org/mongo-driver v1.5.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.32.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.7.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.7.0 // indirect
	go.opentelemetry.io/contrib/samplers/jaegerremote v0.2.0 // indirect
	go.opentelemetry.io/otel/bridge/opentracing v1.6.3 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.opentelemetry.io/proto/otlp v0.16.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20220517005047-85d78b3ac167 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20220524220425-1d687d428aca // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220525015930-6ca3db687a9d // indirect
	google.golang.org/grpc v1.46.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)
