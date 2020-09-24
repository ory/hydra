module github.com/ory/hydra

go 1.14

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
	github.com/Microsoft/go-winio v0.4.11
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5
	github.com/asaskevich/govalidator v0.0.0-20170425121227-4918b99a7cb9
	github.com/aws/aws-sdk-go v1.26.5 // indirect
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/containerd/continuity v0.0.0-20180921161001-7f53d412b9eb
	github.com/coupa/foundation-go v1.2.2
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.0.0+incompatible
	github.com/docker/go-connections v0.3.0
	github.com/docker/go-units v0.3.3
	github.com/fsnotify/fsnotify v1.4.3-0.20170329110642-4da3e2cfbabc
	github.com/go-sql-driver/mysql v1.3.0
	github.com/golang/protobuf v0.0.0-20170920220647-130e6b02ab05
	github.com/gorilla/context v0.0.0-20160226214623-1ea25387ff6f
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v0.0.0-20160922145804-ca9ada445741
	github.com/hashicorp/golang-lru v0.5.0
	github.com/hashicorp/hcl v0.0.0-20171009174708-42e33e2d55a0
	github.com/imdario/mergo v0.0.0-20160216103600-3e95a51e0639
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/jehiah/go-strftime v0.0.0-20140929120216-834e15c05a45
	github.com/jmoiron/sqlx v0.0.0-20180614180643-0dae4fefe7c0
	github.com/julienschmidt/httprouter v1.1.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/lib/pq v0.0.0-20171010183604-30d59eaf0152
	github.com/magiconair/properties v1.7.4-0.20170902060319-8d7837e64d3c
	github.com/mattn/go-sqlite3 v2.0.2+incompatible // indirect
	github.com/meatballhat/negroni-logrus v0.0.0-20170801195057-31067281800f
	github.com/mitchellh/mapstructure v0.0.0-20170523030023-d0303fe80992
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/moul/http2curl v1.0.0
	github.com/oleiade/reflections v1.0.0
	github.com/opencontainers/go-digest v1.0.0-rc1.0.20180430190053-c9281466c8b2
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runc v1.0.0-rc5
	github.com/ory/dockertest v3.3.2+incompatible
	github.com/ory/fosite v0.10.1-0.20191218201927-75c261d8b590
	github.com/ory/graceful v0.1.0
	github.com/ory/herodot v0.1.1
	github.com/ory/ladon v0.8.9
	github.com/ory/pagination v0.0.2-0.20180227110002-05947c3e39e2
	github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/pelletier/go-toml v1.0.2-0.20171001224747-2009e44b6f18
	github.com/pkg/errors v0.8.0
	github.com/pkg/profile v1.2.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/rubenv/sql-migrate v0.0.0-20180704111356-3f452fc0ebeb
	github.com/segmentio/analytics-go v2.0.1-0.20160711225931-bdb0aeca8a99+incompatible
	github.com/segmentio/backo-go v0.0.0-20160424052352-204274ad699c
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v0.0.0-20171008182726-e67d870304c4
	github.com/spf13/cast v1.1.0
	github.com/spf13/cobra v0.0.0-20171010034433-4d6af280c76f
	github.com/spf13/jwalterweatherman v0.0.0-20170901151539-12bd96e66386
	github.com/spf13/pflag v1.0.1-0.20171008183526-a9789e855c76
	github.com/spf13/viper v1.0.1-0.20170929210642-d9cca5ef3303
	github.com/square/go-jose v2.1.3+incompatible
	github.com/stretchr/testify v1.2.2
	github.com/toqueteos/webbrowser v1.0.0
	github.com/urfave/negroni v0.2.0
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c
	golang.org/x/crypto v0.0.0-20170930174604-9419663f5a44
	golang.org/x/net v0.0.0-20180925072008-f04abc6bdfa7
	golang.org/x/oauth2 v0.0.0-20170928010508-bb50c06baba3
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894
	golang.org/x/text v0.1.1-0.20171006144033-825fc78a2fd6
	google.golang.org/appengine v1.0.1-0.20171010223110-07f075729064
	gopkg.in/airbrake/gobrake.v2 v2.0.8
	gopkg.in/alexcesaro/statsd.v2 v2.0.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/gorp.v1 v1.7.1
	gopkg.in/square/go-jose.v2 v2.1.3
	gopkg.in/yaml.v2 v2.0.0-20170812160011-eb3733d160e7
)

replace github.com/ory/fosite => github.com/coupa/fosite v0.10.1-0.20191218201927-75c261d8b590
