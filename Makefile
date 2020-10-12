SHELL=/bin/bash -o pipefail

export GO111MODULE := on
export PATH := .bin:${PATH}

GO_DEPENDENCIES = github.com/ory/go-acc \
				  github.com/sqs/goreturns \
				  github.com/ory/x/tools/listx \
				  github.com/golang/mock/mockgen \
				  github.com/go-swagger/go-swagger/cmd/swagger \
				  golang.org/x/tools/cmd/goimports \
				  github.com/ory/cli \
				  github.com/gobuffalo/packr/v2/packr2 \
				  github.com/go-bindata/go-bindata/go-bindata

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  .bin/$(notdir $1): go.sum go.mod
		GOBIN=$(PWD)/.bin/ go install $1
endef

.bin/golangci-lint: Makefile
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b .bin v1.31.0

$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))

node_modules: package.json
		npm ci

.bin/clidoc:
		go build -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
		clidoc .

.PHONY: lint
lint: .bin/golangci-lint
		golangci-lint run -v ./...

# Runs full test suite including tests where databases are enabled
.PHONY: test-legacy-migrations
test-legacy-migrations: test-resetdb sqlbin
		source scripts/test-env.sh && go test -tags legacy_migration_test sqlite -failfast -timeout=20m ./internal/fizzmigrate
		docker rm -f hydra_test_database_mysql
		docker rm -f hydra_test_database_postgres
		docker rm -f hydra_test_database_cockroach

# Runs full test suite including tests where databases are enabled
.PHONY: test
test: .bin/go-acc
		make test-resetdb
		source scripts/test-env.sh && go-acc ./... -- -failfast -timeout=20m -tags sqlite
		docker rm -f hydra_test_database_mysql
		docker rm -f hydra_test_database_postgres
		docker rm -f hydra_test_database_cockroach

# Resets the test databases
.PHONY: test-resetdb
test-resetdb: node_modules
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
docker: .bin/packr2
		packr2
		GO111MODULE=on GOOS=linux GOARCH=amd64 go build -tags sqlite
		packr2 clean
		docker build -t oryd/hydra:latest .
		docker build -f Dockerfile-alpine -t oryd/hydra:latest-alpine .
		rm hydra

.PHONY: e2e
e2e: node_modules test-resetdb
		source ./scripts/test-env.sh
		./test/e2e/circle-ci.bash memory
		./test/e2e/circle-ci.bash memory-jwt
		./test/e2e/circle-ci.bash postgres
		./test/e2e/circle-ci.bash postgres-jwt
		./test/e2e/circle-ci.bash mysql
		./test/e2e/circle-ci.bash mysql-jwt
		./test/e2e/circle-ci.bash cockroach
		./test/e2e/circle-ci.bash cockroach-jwt

# Runs tests in short mode, without database adapters
.PHONY: quicktest
quicktest:
		go test -failfast -short -tags sqlite ./...

# Formats the code
.PHONY: format
format: .bin/goreturns node_modules
		goreturns -w -local github.com/ory $$(listx .)
		npm run format

# Generates mocks
.PHONY: mocks
mocks: .bin/mockgen
		mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider

# Adds sql files to the binary using go-bindata
.PHONY: sqlbin
sqlbin: .bin/go-bindata
		cd internal/fizzmigrate/client; go-bindata -o sql_migration_files.go -pkg client ./migrations/sql/...
		cd internal/fizzmigrate/consent; go-bindata -o sql_migration_files.go -pkg consent ./migrations/sql/...
		cd internal/fizzmigrate/jwk; go-bindata -o sql_migration_files.go -pkg jwk ./migrations/sql/...
		cd internal/fizzmigrate/oauth2; go-bindata -o sql_migration_files.go -pkg oauth2 ./migrations/sql/...

# Runs all code generators
.PHONY: gen
gen: mocks sqlbin sdk

# Generates the SDKs
.PHONY: sdk
sdk: .bin/cli
		swagger generate spec -m -o ./.schema/api.swagger.json -x internal/httpclient -x gopkg.in/square/go-jose.v2
		cli dev swagger sanitize ./.schema/api.swagger.json
		swagger flatten --with-flatten=remove-unused -o ./.schema/api.swagger.json ./.schema/api.swagger.json
		swagger validate ./.schema/api.swagger.json
		rm -rf internal/httpclient
		mkdir -p internal/httpclient
		swagger generate client -f ./.schema/api.swagger.json -t internal/httpclient -A Ory_Hydra
		make format


.PHONY: install-stable
install-stable: .bin/packr2
		HYDRA_LATEST=$$(git describe --abbrev=0 --tags)
		git checkout $$HYDRA_LATEST
		packr2
		GO111MODULE=on go install \
				-tags sqlite \
				-ldflags "-X github.com/ory/hydra/cmd.Version=$$HYDRA_LATEST -X github.com/ory/hydra/cmd.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/cmd.Commit=`git rev-parse HEAD`" \
				.
		packr2 clean
		git checkout master

.PHONY: install
install: .bin/packr2
		packr2
		GO111MODULE=on go install -tags sqlite .
		packr2 clean
