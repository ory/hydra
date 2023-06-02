SHELL=/bin/bash -o pipefail

export GO111MODULE := on
export PATH := .bin:${PATH}
export PWD := $(shell pwd)

GOLANGCI_LINT_VERSION = 1.46.2

GO_DEPENDENCIES = github.com/ory/go-acc \
				  github.com/golang/mock/mockgen \
				  golang.org/x/tools/cmd/goimports \
				  github.com/go-swagger/go-swagger/cmd/swagger

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  .bin/$(notdir $1): go.sum go.mod
		GOBIN=$(PWD)/.bin/ go install $1
endef

.bin/golangci-lint-$(GOLANGCI_LINT_VERSION):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin v$(GOLANGCI_LINT_VERSION)
	mv .bin/golangci-lint .bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))

node_modules: package-lock.json
	npm ci
	touch node_modules

.PHONY: .bin/yq
.bin/yq:
	go build -o .bin/yq github.com/mikefarah/yq/v4

.bin/clidoc: go.mod
	go build -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
	clidoc .

.bin/licenses: Makefile
	curl https://raw.githubusercontent.com/ory/ci/master/licenses/install | sh

.bin/ory: Makefile
	curl https://raw.githubusercontent.com/ory/meta/master/install.sh | bash -s -- -b .bin ory v0.2.2
	touch .bin/ory

.PHONY: lint
lint: .bin/golangci-lint-$(GOLANGCI_LINT_VERSION)
	.bin/golangci-lint-$(GOLANGCI_LINT_VERSION) run -v ./...

# Runs full test suite including tests where databases are enabled
.PHONY: test
test: .bin/go-acc
	make test-resetdb
	source scripts/test-env.sh && go-acc ./... -- -failfast -timeout=20m -tags sqlite,json1
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
	docker run --rm --name hydra_test_database_mysql -p 3444:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0
	docker run --rm --name hydra_test_database_postgres -p 3445:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=postgres -e PGUSER=postgres -e POSTGRES_DB=postgres -d postgres:15.2
	docker run --rm --name hydra_test_database_cockroach -p 3446:26257 -p 8082:8080 -d cockroachdb/cockroach:v22.2.8 start-single-node --insecure

# Build local docker images
.PHONY: docker
docker:
	docker build -f .docker/Dockerfile-build -t oryd/hydra:latest-sqlite .

.PHONY: bench
bench:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -c ./oauth2/ -o hydrabench
	# docker build -f .docker/Dockerfile-hydrabench --platform linux/amd64 -t orysh.azurecr.io/hydra:bench .
	# docker push orysh.azurecr.io/hydra:bench
	docker build -f .docker/Dockerfile-hydrabench --platform linux/amd64 -t eu.gcr.io/oasis-test-infra/hydra:cookiebench .
	docker push eu.gcr.io/oasis-test-infra/hydra:cookiebench

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
	go test -failfast -short -tags sqlite,json1 ./...

.PHONY: quicktest-hsm
quicktest-hsm:
	docker build --progress=plain -f .docker/Dockerfile-hsm --target test-hsm .

.PHONY: refresh
refresh:
	UPDATE_SNAPSHOTS=true go test -failfast -short -tags sqlite,json1 ./...

authors:  # updates the AUTHORS file
	curl https://raw.githubusercontent.com/ory/ci/master/authors/authors.sh | env PRODUCT="Ory Hydra" bash

# Formats the code
.PHONY: format
format: .bin/goimports .bin/ory node_modules
	.bin/ory dev headers copyright --type=open-source --exclude=internal/httpclient
	.bin/goimports -w --local github.com/ory .
	npm exec -- prettier --write .

# Generates mocks
.PHONY: mocks
mocks: .bin/mockgen
	mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider
	mockgen -package jwk_test -destination jwk/registry_mock_test.go -source=jwk/registry.go
	go generate ./...

# Generates the SDKs
.PHONY: sdk
sdk: .bin/swagger .bin/ory node_modules
	swagger generate spec -m -o spec/swagger.json \
		-c github.com/ory/hydra/v2/client \
		-c github.com/ory/hydra/v2/consent \
		-c github.com/ory/hydra/v2/health \
		-c github.com/ory/hydra/v2/jwk \
		-c github.com/ory/hydra/v2/oauth2 \
		-c github.com/ory/hydra/v2/x \
		-c github.com/ory/x/healthx \
		-c github.com/ory/x/openapix \
		-c github.com/ory/x/pagination \
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
				-p file://.schema/openapi/patches/nulls.yaml \
				-p file://.schema/openapi/patches/common.yaml \
				-p file://.schema/openapi/patches/security.yaml \
				spec/swagger.json spec/api.json

	rm -rf "internal/httpclient"
	npm run openapi-generator-cli -- generate -i "spec/api.json" \
		-g go \
		-o "internal/httpclient" \
		--git-user-id ory \
		--git-repo-id hydra-client-go \
		--git-host github.com
	(cd internal/httpclient && go mod edit -module github.com/ory/hydra-client-go/v2)

	make format

MIGRATIONS_SRC_DIR = ./persistence/sql/src/
MIGRATIONS_DST_DIR = ./persistence/sql/migrations/
MIGRATION_NAMES=$(shell find $(MIGRATIONS_SRC_DIR) -maxdepth 1 -mindepth 1 -type d -print0 | xargs -0 -I{} basename {})
MIGRATION_TARGETS=$(addprefix $(MIGRATIONS_DST_DIR), $(MIGRATION_NAMES))
.PHONY: $(MIGRATION_TARGETS)
$(MIGRATION_TARGETS): $(MIGRATIONS_DST_DIR)%:
	go run . migrate gen $(MIGRATIONS_SRC_DIR)$* $(MIGRATIONS_DST_DIR)
	./scripts/db-placeholders.sh

MIGRATION_CLEAN_TARGETS=$(addsuffix -clean, $(MIGRATION_TARGETS))
.PHONY: $(MIGRATION_CLEAN_TARGETS)
$(MIGRATION_CLEAN_TARGETS): $(MIGRATIONS_DST_DIR)%:
	find $(MIGRATIONS_DST_DIR) -type f -name $$(echo "$*" | cut -c1-14)* -delete


.PHONY: $(MIGRATIONS_DST_DIR:%/=%)
$(MIGRATIONS_DST_DIR:%/=%): $(MIGRATION_TARGETS)

.PHONY: $(MIGRATIONS_DST_DIR:%/=%-clean)
$(MIGRATIONS_DST_DIR:%/=%-clean): $(MIGRATION_CLEAN_TARGETS)

.PHONY: install-stable
install-stable:
	HYDRA_LATEST=$$(git describe --abbrev=0 --tags)
	git checkout $$HYDRA_LATEST
	GO111MODULE=on go install \
		-tags sqlite,json1 \
		-ldflags "-X github.com/ory/hydra/v2/driver/config.Version=$$HYDRA_LATEST -X github.com/ory/hydra/v2/driver/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/v2/driver/config.Commit=`git rev-parse HEAD`" \
		.
	git checkout master

.PHONY: install
install:
	GO111MODULE=on go install -tags sqlite,json1 .

.PHONY: post-release
post-release: .bin/yq
	yq e '.services.hydra.image = "oryd/hydra:'$$DOCKER_TAG'"' -i quickstart.yml
	yq e '.services.hydra-migrate.image = "oryd/hydra:'$$DOCKER_TAG'"' -i quickstart.yml
	yq e '.services.consent.image = "oryd/hydra-login-consent-node:'$$DOCKER_TAG'"' -i quickstart.yml

generate: .bin/mockgen
	go generate ./...

licenses: .bin/licenses node_modules  # checks open-source licenses
	.bin/licenses
