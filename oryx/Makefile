SHELL=/bin/bash -o pipefail

export PATH := .bin:${PATH}

.bin/ory: Makefile
	curl --retry 7 --retry-connrefused https://raw.githubusercontent.com/ory/meta/master/install.sh | bash -s -- -b .bin ory v0.2.2
	touch .bin/ory

.PHONY: format
format: .bin/ory node_modules
	.bin/ory dev headers copyright --type=open-source --exclude=clidoc/ --exclude=hasherx/mocks_pkdbf2_test.go --exclude=josex/ --exclude=hasherx/ --exclude=jsonnetsecure/jsonnet.go
	go tool goimports -w -local github.com/ory .
	npm exec -- prettier --write .

.bin/golangci-lint: Makefile
	curl --retry 7 --retry-connrefused -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin v2.4.0

.bin/licenses: Makefile
	curl --retry 7 --retry-connrefused https://raw.githubusercontent.com/ory/ci/master/licenses/install | sh

licenses: .bin/licenses node_modules  # checks open-source licenses
	.bin/licenses

.PHONY: test
test:
	make resetdb
	export TEST_DATABASE_POSTGRESQL=postgres://postgres:secret@127.0.0.1:3445/hydra?sslmode=disable; export TEST_DATABASE_COCKROACHDB=cockroach://root@127.0.0.1:3446/defaultdb?sslmode=disable; export TEST_DATABASE_MYSQL='mysql://root:secret@tcp(127.0.0.1:3444)/mysql?parseTime=true&multiStatements=true'; go test -count=1 -tags sqlite ./...

.PHONY: resetdb
resetdb:
	docker kill hydra_test_database_mysql || true
	docker kill hydra_test_database_postgres || true
	docker kill hydra_test_database_cockroach || true
	docker rm -f hydra_test_database_mysql || true
	docker rm -f hydra_test_database_postgres || true
	docker rm -f hydra_test_database_cockroach || true
	docker run --rm --name hydra_test_database_mysql -p 3444:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0
	docker run --rm --name hydra_test_database_postgres -p 3445:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=hydra -d postgres:11.8
	docker run --rm --name hydra_test_database_cockroach -p 3446:26257 -d cockroachdb/cockroach:latest-v25.3 start-single-node --insecure

.PHONY: lint
lint: .bin/golangci-lint
	GO111MODULE=on .bin/golangci-lint run -v ./...

.PHONY: migrations-render
migrations-render: .bin/ory
	ory dev pop migration render networkx/migrations/templates networkx/migrations/sql

.PHONY: migrations-render-replace
migrations-render-replace: .bin/ory
	ory dev pop migration render -r networkx/migrations/templates networkx/migrations/sql

.PHONY: mocks
mocks:
	go tool mockgen -package hasherx_test -destination hasherx/mocks_argon2_test.go github.com/ory/x/hasherx Argon2Configurator
	go tool mockgen -package hasherx_test -destination hasherx/mocks_bcrypt_test.go github.com/ory/x/hasherx BCryptConfigurator
	go tool mockgen -package hasherx_test -destination hasherx/mocks_pkdbf2_test.go github.com/ory/x/hasherx PBKDF2Configurator

node_modules: package-lock.json
	npm ci
	touch node_modules
