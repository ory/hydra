init:
		go get -u \
			github.com/ory/x/tools/listx \
			github.com/sqs/goreturns \
			github.com/golang/mock/mockgen \
			github.com/go-swagger/go-swagger/cmd/swagger \
			github.com/go-bindata/go-bindata/... \
			golang.org/x/tools/cmd/goimports \
			github.com/gobuffalo/packr/packr

format:
		goreturns -w -local github.com/ory $$(listx .)
		# goimports -w -local github.com/ory $$(listx .)

gen-mocks:
		mockgen -package oauth2_test -destination oauth2/oauth2_provider_mock_test.go github.com/ory/fosite OAuth2Provider

gen-sql:
		cd client; go-bindata -o sql_migration_files.go -pkg client ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests
		cd consent; go-bindata -o sql_migration_files.go -pkg consent ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests
		cd jwk; go-bindata -o sql_migration_files.go -pkg jwk ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests
		cd oauth2; go-bindata -o sql_migration_files.go -pkg oauth2 ./migrations/sql/shared ./migrations/sql/mysql ./migrations/sql/postgres ./migrations/sql/tests

gen: gen-mocks gen-sql gen-sdk

gen-sdk:
		swagger generate spec -m -o ./docs/api.swagger.json
		swagger validate ./docs/api.swagger.json

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
