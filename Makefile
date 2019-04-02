SHELL=/bin/bash -o pipefail

# Runs full test suite including tests where databases are enabled
.PHONY: test
test:
		docker kill hydra_test_database_mysql || true
		docker kill hydra_test_database_postgres || true
		docker rm -f hydra_test_database_mysql || true
		docker rm -f hydra_test_database_postgres || true
		docker run --rm --name hydra_test_database_mysql -p 3444:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:5.7
		docker run --rm --name hydra_test_database_postgres -p 3445:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=hydra -d postgres:9.6
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
.PHONY: test-short
test-short:
		go test -failfast -short ./...

# Formats the code
.PHONY: format
format:
		goreturns -w -local github.com/ory $$(listx .)

# Generates mocks
.PHONY: mocks
mocks:
		mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider
		mockgen -mock_names Provider=MockConfiguration -package internal -destination internal/configuration_provider_mock.go github.com/ory/hydra/driver/configuration Provider
		mockgen -mock_names Registry=MockRegistry -package internal -destination internal/registry_mock.go github.com/ory/hydra/driver Registry


# Adds sql files to the binary using go-bindata
.PHONY: sqlbin
sqlbin:
		cd client; go-bindata -o sql_migration_files.go -pkg client ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests
		cd consent; go-bindata -o sql_migration_files.go -pkg consent ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests
		cd jwk; go-bindata -o sql_migration_files.go -pkg jwk ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests
		cd oauth2; go-bindata -o sql_migration_files.go -pkg oauth2 ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests

# Runs all code generators
.PHONY: gen
gen: mocks sqlbin sdks

# Generates the SDKs
.PHONY: sdks
sdks:
		GO111MODULE=on go mod tidy
		GO111MODULE=on go mod vendor
		GO111MODULE=off swagger generate spec -m -o ./docs/api.swagger.json
		GO111MODULE=off swagger validate ./docs/api.swagger.json

		rm -rf ./sdk/go/hydra/swagger
		rm -rf ./sdk/js/swagger
		rm -rf ./sdk/php/swagger
		rm -rf ./sdk/java

		# swagger generate client -f ./docs/api.swagger.json -t sdk/go
		java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l go -o ./sdk/go/hydra/swagger
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

		git checkout HEAD -- sdk/go/hydra/swagger/configuration.go
		git checkout HEAD -- sdk/go/hydra/swagger/api_client.go

		rm -f ./sdk/js/swagger/package.json
		rm -rf ./sdk/js/swagger/test
		rm -f ./sdk/php/swagger/composer.json ./sdk/php/swagger/phpunit.xml.dist
		rm -rf ./sdk/php/swagger/test
		rm -rf ./vendor
