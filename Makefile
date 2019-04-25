SHELL=/bin/bash -o pipefail

# Runs full test suite including tests where databases are enabled
.PHONY: test
test:
		make test-resetdb
		make sqlbin
		TEST_DATABASE_MYSQL='root:secret@(127.0.0.1:3444)/mysql?parseTime=true' \
			TEST_DATABASE_POSTGRESQL='postgres://postgres:secret@127.0.0.1:3445/hydra?sslmode=disable' \
			go-acc ./... -- -failfast -timeout=20m
		docker rm -f hydra_test_database_mysql
		docker rm -f hydra_test_database_postgres

# Resets the test databases
.PHONY: test-resetdb
test-resetdb:
		docker kill hydra_test_database_mysql || true
		docker kill hydra_test_database_postgres || true
		docker rm -f hydra_test_database_mysql || true
		docker rm -f hydra_test_database_postgres || true
		docker run --rm --name hydra_test_database_mysql -p 3444:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:5.7
		docker run --rm --name hydra_test_database_postgres -p 3445:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=hydra -d postgres:9.6

# Runs tests in short mode, without database adapters
.PHONY: docker
docker:
		make sqlbin
		GO111MODULE=on GOOS=linux GOARCH=amd64 go build
		docker build -t oryd/hydra:latest .
		rm hydra

.PHONY: e2e
e2e:
		make e2e-memory
		make e2e-memory-jwt
		make e2e-postgres
		make e2e-postgres-jwt
		make e2e-mysql
		make e2e-mysql-jwt

.PHONY: e2e-memory
e2e-memory:
		make e2e-prepare-memory
		make e2e-waiton
		npm run test

.PHONY: e2e-memory-jwt
e2e-memory-jwt:
		make e2e-prepare-memory-jwt
		make e2e-waiton
		npm run test

.PHONY: e2e-postgres
e2e-postgres:
		make e2e-prepare-postgres
		make e2e-waiton
		npm run test

.PHONY: e2e-postgres-jwt
e2e-postgres-jwt:
		make e2e-prepare-postgres-jwt
		make e2e-waiton
		npm run test

.PHONY: e2e-mysql
e2e-mysql:
		make e2e-prepare-mysql
		make e2e-waiton
		npm run test

.PHONY: e2e-mysql-jwt
e2e-mysql-jwt:
		make e2e-prepare-mysql-jwt
		make e2e-waiton
		npm run test

.PHONY: e2e-plugin
e2e-plugin:
		make e2e-prepare-plugin
		make e2e-waiton
		npm run test

.PHONY: e2e-plugin-jwt
e2e-plugin-jwt:
		make e2e-prepare-plugin-jwt
		make e2e-waiton
		npm run test

.PHONY: e2e-waiton
e2e-waiton:
		npm run wait-on http-get://localhost:5000/health/ready
		npm run wait-on http-get://localhost:5001/health/ready
		npm run wait-on http-get://localhost:5002/
		npm run wait-on http-get://localhost:5003/oauth2/callback

.PHONY: e2e-prepare-memory
e2e-prepare-memory:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			up --build -d

.PHONY: e2e-prepare-memory-jwt
e2e-prepare-memory-jwt:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.jwt.yml \
			up --build -d

.PHONY: e2e-prepare-postgres
e2e-prepare-postgres:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.postgres.yml \
			up --build -d

.PHONY: e2e-prepare-postgres-jwt
e2e-prepare-postgres-jwt:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.postgres.yml \
			-f ./test/e2e/docker-compose.jwt.yml \
			up --build -d

.PHONY: e2e-prepare-mysql
e2e-prepare-mysql:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.mysql.yml \
			up --build -d
			
.PHONY: e2e-prepare-mysql-jwt
e2e-prepare-mysql-jwt:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.mysql.yml \
			-f ./test/e2e/docker-compose.jwt.yml \
			up --build -d

.PHONY: e2e-prepare-plugin
e2e-prepare-plugin:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.plugin.yml \
			up --build -d
			
.PHONY: e2e-prepare-plugin-jwt
e2e-prepare-plugin-jwt:
		make e2e-prepare-reset
		docker-compose \
			-f ./test/e2e/docker-compose.yml \
			-f ./test/e2e/docker-compose.plugin.yml \
			-f ./test/e2e/docker-compose.jwt.yml \
			up --build -d

.PHONY: e2e-prepare-reset
e2e-prepare-reset:
		docker-compose \
			-f ./test/e2e/docker-compose.jwt.yml \
			-f ./test/e2e/docker-compose.mysql.yml \
			-f ./test/e2e/docker-compose.plugin.yml \
			-f ./test/e2e/docker-compose.postgres.yml \
			-f ./test/e2e/docker-compose.yml \
			kill
		docker-compose \
			-f ./test/e2e/docker-compose.jwt.yml \
			-f ./test/e2e/docker-compose.mysql.yml \
			-f ./test/e2e/docker-compose.plugin.yml \
			-f ./test/e2e/docker-compose.postgres.yml \
			-f ./test/e2e/docker-compose.yml \
			rm -f
		make docker

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
		cd client; go-bindata -o sql_migration_files.go -pkg client ./migrations/sql/...
		cd consent; go-bindata -o sql_migration_files.go -pkg consent ./migrations/sql/...
		cd jwk; go-bindata -o sql_migration_files.go -pkg jwk ./migrations/sql/...
		cd oauth2; go-bindata -o sql_migration_files.go -pkg oauth2 ./migrations/sql/...

# Runs all code generators
.PHONY: gen
gen: mocks sqlbin sdk

# Generates the SDKs
.PHONY: sdk
sdk:
		rm -rf ./vendor/
		GO111MODULE=on go mod tidy
		GO111MODULE=on go mod vendor
		GO111MODULE=off swagger generate spec -m -o ./docs/api.swagger.json
		GO111MODULE=off swagger validate ./docs/api.swagger.json

		rm -rf ./sdk/go/hydra
		rm -rf ./sdk/js/swagger
		rm -rf ./sdk/php/swagger
		rm -rf ./sdk/java

		mkdir ./sdk/go/hydra

		GO111MODULE=off swagger generate client -f ./docs/api.swagger.json -t sdk/go/hydra -A Ory_Hydra
		java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l javascript -o ./sdk/js/swagger
		java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l php -o sdk/php/ \
			--invoker-package Hydra\\SDK --git-repo-id swagger --git-user-id ory --additional-properties "packagePath=swagger,description=Client for Hydra"
		java -DapiTests=false -DmodelTests=false -jar scripts/swagger-codegen-cli-2.2.3.jar generate \
			--input-spec ./docs/api.swagger.json \
			--lang java \
			--library resttemplate \
			--group-id com.github.ory \
			--artifact-id hydra-client-resttemplate \
			--invoker-package com.github.ory.hydra \
			--api-package com.github.ory.hydra.api \
			--model-package com.github.ory.hydra.model \
			--output ./sdk/java/hydra-client-resttemplate

		cd sdk/go; goreturns -w -i -local github.com/ory $$(listx .)

		rm -f ./sdk/js/swagger/package.json
		rm -rf ./sdk/js/swagger/test
		rm -f ./sdk/php/swagger/composer.json ./sdk/php/swagger/phpunit.xml.dist
		rm -rf ./sdk/php/swagger/test
		rm -rf ./vendor

.PHONY: install-stable
install-stable:
		HYDRA_LATEST=$$(git describe --abbrev=0 --tags)
		git checkout $$HYDRA_LATEST
		GO111MODULE=on go install \
				-ldflags "-X github.com/ory/hydra/cmd.Version=$$HYDRA_LATEST -X github.com/ory/hydra/cmd.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/cmd.Commit=`git rev-parse HEAD`" \
				.
		git checkout master

.PHONY: install
install:
		GO111MODULE=on go install .
