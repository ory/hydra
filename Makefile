SHELL=/bin/bash -o pipefail

export GO111MODULE := on
export PATH := .bin:${PATH}
export PWD := $(shell pwd)

GO_DEPENDENCIES = github.com/ory/go-acc \
				  golang.org/x/tools/cmd/goimports \
				  github.com/golang/mock/mockgen \
				  github.com/go-swagger/go-swagger/cmd/swagger \
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

docs/node_modules: docs/package.json
		cd docs; npm ci

.bin/clidoc: go.mod
		go build -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
		clidoc .

.bin/ory: Makefile
		bash <(curl https://raw.githubusercontent.com/ory/cli/master/install.sh) -b .bin v0.0.53
		touch -a -m .bin/ory

.PHONY: lint
lint: .bin/golangci-lint
		golangci-lint run -v ./...

# Runs full test suite including tests where databases are enabled
.PHONY: test-legacy-migrations
test-legacy-migrations: test-resetdb .bin/go-bindata
		cd internal/fizzmigrate/client; go-bindata -o sql_migration_files.go -pkg client ./migrations/sql/...
		cd internal/fizzmigrate/consent; go-bindata -o sql_migration_files.go -pkg consent ./migrations/sql/...
		cd internal/fizzmigrate/jwk; go-bindata -o sql_migration_files.go -pkg jwk ./migrations/sql/...
		cd internal/fizzmigrate/oauth2; go-bindata -o sql_migration_files.go -pkg oauth2 ./migrations/sql/...
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
		docker run --rm --name hydra_test_database_cockroach -p 3446:26257 -d cockroachdb/cockroach:v20.2.6 start-single-node --insecure

# Runs tests in short mode, without database adapters
.PHONY: docker
docker:
		docker build -f .docker/Dockerfile-build -t oryd/hydra:latest-sqlite .

.PHONY: e2e
e2e: node_modules test-resetdb
		source ./scripts/test-env.sh
		for db in memory postgres mysql cockroach; do \
			./test/e2e/circle-ci.bash "$${db}"; \
			./test/e2e/circle-ci.bash "$${db}" --jwt; \
			./test/e2e/circle-ci.bash "$${db}" --grant_jwt_client_auth_optional --grant_jwt_jti_optional --grant_jwt_iat_optional; \
		done

# Runs tests in short mode, without database adapters
.PHONY: quicktest
quicktest:
		go test -failfast -short -tags sqlite ./...

# Formats the code
.PHONY: format
format: .bin/goimports node_modules docs/node_modules
		goimports -w --local github.com/ory .
		npm run format
		cd docs; npm run format

# Generates mocks
.PHONY: mocks
mocks: .bin/mockgen
		mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider

# Generates the SDKs
.PHONY: sdk
sdk: .bin/ory
		swagger generate spec -m -o ./spec/api.json -x internal/httpclient -x gopkg.in/square/go-jose.v2
		ory dev swagger sanitize ./spec/api.json
		swagger flatten --with-flatten=remove-unused -o ./spec/api.json ./spec/api.json
		swagger validate ./spec/api.json
		rm -rf internal/httpclient
		mkdir -p internal/httpclient
		swagger generate client -f ./spec/api.json -t internal/httpclient -A Ory_Hydra
		make format

.PHONY: install-stable
install-stable:
		HYDRA_LATEST=$$(git describe --abbrev=0 --tags)
		git checkout $$HYDRA_LATEST
		GO111MODULE=on go install \
				-tags sqlite \
				-ldflags "-X github.com/ory/hydra/driver/config.Version=$$HYDRA_LATEST -X github.com/ory/hydra/driver/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/driver/config.Commit=`git rev-parse HEAD`" \
				.
		git checkout master

.PHONY: install
install:
		GO111MODULE=on go install -tags sqlite .
