SHELL=/bin/bash -o pipefail

export PATH := $(pwd)/.bin:${PATH}

.PHONY: deps
deps:
ifneq ("$(shell cat Makefile | md5))","$(shell cat .bin/.lock)")
		npm ci
		go build -o .bin/go-acc github.com/ory/go-acc
		go build -o .bin/goreturns github.com/sqs/goreturns
		go build -o .bin/listx github.com/ory/x/tools/listx
		go build -o .bin/mockgen github.com/golang/mock/mockgen
		go build -o .bin/swagger github.com/go-swagger/go-swagger/cmd/swagger
		go build -o .bin/goimports golang.org/x/tools/cmd/goimports
		go build -o .bin/swagutil github.com/ory/sdk/swagutil
		go build -o .bin/packr2 github.com/gobuffalo/packr/v2/packr2
		go build -o .bin/go-bindata github.com/go-bindata/go-bindata/go-bindata
		echo "v0" > .bin/.lock
		echo "$$(cat Makefile | md5)" > .bin/.lock
endif

# Runs full test suite including tests where databases are enabled
.PHONY: test-legacy-migrations
test-legacy-migrations:
		make test-resetdb
		make sqlbin
		source scripts/test-env.sh && go test -tags legacy_migration_test -failfast -timeout=20m ./internal/fizzmigrate
		docker rm -f hydra_test_database_mysql
		docker rm -f hydra_test_database_postgres
		docker rm -f hydra_test_database_cockroach

# Runs full test suite including tests where databases are enabled
.PHONY: test
test:
		make test-resetdb
		source scripts/test-env.sh && go-acc ./... -- -failfast -timeout=20m
		docker rm -f hydra_test_database_mysql
		docker rm -f hydra_test_database_postgres
		docker rm -f hydra_test_database_cockroach

# Resets the test databases
.PHONY: test-resetdb
test-resetdb:
		docker kill hydra_test_database_mysql || true
		docker kill hydra_test_database_postgres || true
		docker kill hydra_test_database_cockroach || true
		docker rm -f hydra_test_database_mysql || true
		docker rm -f hydra_test_database_postgres || true
		docker rm -f hydra_test_database_cockroach || true
		docker run --rm --name hydra_test_database_mysql -p 3444:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:5.7
		docker run --rm --name hydra_test_database_postgres -p 3445:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=postgres -d postgres:9.6
		docker run --rm --name hydra_test_database_cockroach -p 3446:26257 -d cockroachdb/cockroach:v2.1.6 start --insecure

# Runs tests in short mode, without database adapters
.PHONY: docker
docker:
		packr2
		CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build
		packr2 clean
		docker build -t oryd/hydra:latest .
		docker build -f Dockerfile-alpine -t oryd/hydra:latest-alpine .
		rm hydra

.PHONY: e2e
e2e:
		make test-resetdb
		source ./scripts/test-env.sh
		./test/e2e/circle-ci.bash memory
		./test/e2e/circle-ci.bash memory-jwt
		./test/e2e/circle-ci.bash postgres
		./test/e2e/circle-ci.bash postgres-jwt
		./test/e2e/circle-ci.bash mysql
		./test/e2e/circle-ci.bash mysql-jwt
		./test/e2e/circle-ci.bash cockroach
		./test/e2e/circle-ci.bash cockroach-jwt
		./test/e2e/circle-ci.bash plugin
		./test/e2e/circle-ci.bash plugin-jwt

# Runs tests in short mode, without database adapters
.PHONY: quicktest
quicktest:
		go test -failfast -short ./...

# Formats the code
.PHONY: format
format:
		goreturns -w -local github.com/ory $$(listx .)
		npm run format

# Generates mocks
.PHONY: mocks
mocks:
		mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider

# Adds sql files to the binary using go-bindata
.PHONY: sqlbin
sqlbin:
		cd internal/fizzmigrate/client; go-bindata -o sql_migration_files.go -pkg client ./migrations/sql/...
		cd internal/fizzmigrate/consent; go-bindata -o sql_migration_files.go -pkg consent ./migrations/sql/...
		cd internal/fizzmigrate/jwk; go-bindata -o sql_migration_files.go -pkg jwk ./migrations/sql/...
		cd internal/fizzmigrate/oauth2; go-bindata -o sql_migration_files.go -pkg oauth2 ./migrations/sql/...

# Runs all code generators
.PHONY: gen
gen: mocks sqlbin sdk

# Generates the SDKs
.PHONY: sdk
sdk:
		swagger generate spec -m -o ./.schema/api.swagger.json -x internal/httpclient,gopkg.in/square/go-jose.v2
		swagutil sanitize ./.schema/api.swagger.json
		swagger flatten --with-flatten=remove-unused -o ./.schema/api.swagger.json ./.schema/api.swagger.json
		swagger validate ./.schema/api.swagger.json
		rm -rf internal/httpclient
		mkdir -p internal/httpclient
		swagger generate client -f ./.schema/api.swagger.json -t internal/httpclient -A Ory_Hydra
		make format


.PHONY: install-stable
install-stable:
		HYDRA_LATEST=$$(git describe --abbrev=0 --tags)
		git checkout $$HYDRA_LATEST
		packr2
		GO111MODULE=on go install \
				-ldflags "-X github.com/ory/hydra/cmd.Version=$$HYDRA_LATEST -X github.com/ory/hydra/cmd.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/cmd.Commit=`git rev-parse HEAD`" \
				.
		packr2 clean
		git checkout master

.PHONY: install
install:
		packr2
		GO111MODULE=on go install .
		packr2 clean

.PHONY: init
init:
		GO111MODULE=on go get .
		GO111MODULE=on go install github.com/ory/go-acc
