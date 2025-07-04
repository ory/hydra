// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package modx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const stub = `module github.com/ory/x

// remove once https://github.com/seatgeek/logrus-gelf-formatter/pull/5 is merged
replace github.com/seatgeek/logrus-gelf-formatter => github.com/zepatrik/logrus-gelf-formatter v0.0.0-20210305135027-b8b3731dba10

require (
	github.com/DataDog/datadog-go v4.0.0+incompatible // indirect
	github.com/bmatcuk/doublestar/v2 v2.0.3
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/dgraph-io/ristretto v0.0.2
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20201201034508-7d75c1d40d88+incompatible
	github.com/fatih/structs v1.1.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghodss/yaml v1.0.0
	github.com/go-bindata/go-bindata v3.1.1+incompatible
	github.com/go-openapi/errors v0.20.0 // indirect
	github.com/go-openapi/runtime v0.19.26
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/fizz v1.10.0
	github.com/gobuffalo/httptest v1.0.2
	github.com/gobuffalo/packr v1.22.0
	github.com/ory/pop/v5 v5.3.1
	github.com/golang/mock v1.3.1
	github.com/google/go-jsonnet v0.16.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-retryablehttp v0.6.8
	github.com/inhies/go-bytesize v0.0.0-20201103132853-d0aed0d254f8
	github.com/jackc/pgconn v1.6.0
	github.com/jackc/pgx/v4 v4.6.0
	github.com/jandelgado/gcov2lcov v1.0.4-0.20210120124023-b83752c6dc08
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/knadh/koanf v0.14.1-0.20201201075439-e0853799f9ec
	github.com/lib/pq v1.3.0
	github.com/markbates/pkger v0.17.1
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/ory/analytics-go/v5 v5.0.0
	github.com/ory/dockertest/v3 v3.6.3
	github.com/ory/go-acc v0.2.6
	github.com/ory/herodot v0.9.2
	github.com/ory/jsonschema/v3 v3.0.1
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.8.0
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.2.1
	github.com/rs/cors v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20190212093014-1007f53448d7
	github.com/seatgeek/logrus-gelf-formatter v0.0.0-20210219220335-367fa274be2c
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cast v1.3.2-0.20200723214538-8d17101741c8
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/go-jose/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.3.2
	github.com/tidwall/sjson v1.0.4
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/urfave/negroni v1.0.0
	go.elastic.co/apm v1.8.0
	go.elastic.co/apm/module/apmot v1.8.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gonum.org/v1/plot v0.0.0-20200111075622-4abb28f724d5
	google.golang.org/grpc v1.36.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.27.0
	gopkg.in/square/go-jose.v2 v2.2.2
)

go 1.16
`

func TestVersion(t *testing.T) {
	for _, tc := range [][]string{
		{"google.golang.org/grpc", "v1.36.0"},
		{"golang.org/x/crypto", "v0.0.0-20200510223506-06a226fb4e37"},
	} {

		v, err := FindVersion([]byte(stub), tc[0])
		require.NoError(t, err)
		assert.Equal(t, tc[1], v)

	}

	_, err := FindVersion([]byte(stub), "notgithub.com/idonot/exist")
	require.Error(t, err)
}
