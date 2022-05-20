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
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b .bin v1.44.0

$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))

node_modules: package.json
		npm ci

.PHONY: .bin/yq
.bin/yq:
		go build -o .bin/yq github.com/mikefarah/yq/v4

.bin/clidoc: go.mod
		go build -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
		clidoc .

.bin/ory: Makefile
		bash <(curl https://raw.githubusercontent.com/ory/meta/master/install.sh) -d -b .bin ory v0.1.22
		touch -a -m .bin/ory

.PHONY: lint
lint: .bin/golangci-lint
		golangci-lint run -v ./...

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
		docker run --rm --name hydra_test_database_mysql --platform linux/amd64 -p 3444:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:5.7
		docker run --rm --name hydra_test_database_postgres -p 3445:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=postgres -d postgres:9.6
		docker run --rm --name hydra_test_database_cockroach -p 3446:26257 -d cockroachdb/cockroach:v20.2.6 start-single-node --insecure

# Build local docker images
.PHONY: docker
docker:
		docker build -f .docker/Dockerfile-build -t oryd/hydra:latest-sqlite .

.PHONY: e2e
e2e: node_modules test-resetdb
		source ./scripts/test-env.sh
		for db in memory postgres mysql cockroach; do \
			./test/e2e/circle-ci.bash "$${db}"; \
			./test/e2e/circle-ci.bash "$${db}" --jwt; \
		done

# Runs tests in short mode, without database adapters
.PHONY: quicktest
quicktest:
		go test -failfast -short -tags sqlite ./...

.PHONY: quicktest-hsm
quicktest-hsm:
		docker build --progress=plain -f .docker/Dockerfile-hsm --target test-hsm .

# Formats the code
.PHONY: format
format: .bin/goimports node_modules
		goimports -w --local github.com/ory .
		npm run format

# Generates mocks
.PHONY: mocks
mocks: .bin/mockgen
		mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider

# Generates the SDKs
.PHONY: sdk
sdk: .bin/swagger .bin/ory node_modules
		swagger generate spec -m -o spec/swagger.json \
			-c github.com/ory/hydra/client \
			-c github.com/ory/hydra/consent \
			-c github.com/ory/hydra/health \
			-c github.com/ory/hydra/jwk \
			-c github.com/ory/hydra/oauth2 \
			-c github.com/ory/hydra/x \
			-c github.com/ory/x/healthx \
			-c github.com/ory/herodot
		ory dev swagger sanitize ./spec/swagger.json
		swagger validate ./spec/swagger.json
		CIRCLE_PROJECT_USERNAME=ory CIRCLE_PROJECT_REPONAME=hydra \
				ory dev openapi migrate \
					--health-path-tags metadata \
					-p https://raw.githubusercontent.com/ory/x/master/healthx/openapi/patch.yaml \
					-p file://.schema/openapi/patches/meta.yaml \
					-p file://.schema/openapi/patches/health.yaml \
					-p file://.schema/openapi/patches/oauth2.yaml \
					spec/swagger.json spec/api.json

		rm -rf internal/httpclient internal/httpclient-next
		mkdir -p internal/httpclient internal/httpclient-next
		swagger generate client -f ./spec/swagger.json -t internal/httpclient -A Ory_Hydra

		npm run openapi-generator-cli -- generate -i "spec/api.json" \
				-g go \
				-o "internal/httpclient-next" \
				--git-user-id ory \
				--git-repo-id hydra-client-go \
				--git-host github.com \
				-t .schema/openapi/templates/go \
				-c .schema/openapi/gen.go.yml

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

.PHONY: contributors
contributors:
		printf '# contributors generated by `make contributors`\n\n' > ./CONTRIBUTORS
		git log --format="%aN <%aE>" | sort | uniq | grep -v '^dependabot\[bot\]' >> ./CONTRIBUTORS

.PHONY: post-release
post-release: .bin/yq
		cat quickstart.yml | yq '.services.hydra.image = "oryd/hydra:'$$DOCKER_TAG'"' | sponge quickstart.yml
		cat quickstart.yml | yq '.services.hydra-migrate.image = "oryd/hydra:'$$DOCKER_TAG'"' | sponge quickstart.yml
		cat quickstart.yml | yq '.services.consent.image = "oryd/hydra-login-consent-node:'$$DOCKER_TAG'"' | sponge quickstart.yml
