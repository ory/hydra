export PATH 					:= .bin:${PATH}

format: .bin/goimports .bin/ory node_modules  # formats the source code
	.bin/ory dev headers copyright --type=open-source
	.bin/goimports -w .
	npm exec -- prettier --write .

help:
	@cat Makefile | grep '^[^ ]*:' | grep -v '^\.bin/' | grep -v '.SILENT:' | grep -v '^node_modules:' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

licenses: .bin/licenses node_modules  # checks open-source licenses
	.bin/licenses

test:  # runs all tests
	go test ./...

.bin/goimports: Makefile
	GOBIN=$(shell pwd)/.bin go install golang.org/x/tools/cmd/goimports@latest

.bin/licenses: Makefile
	curl https://raw.githubusercontent.com/ory/ci/master/licenses/install | sh

.bin/mockgen:
	go build -o .bin/mockgen go.uber.org/mock/mockgen

.bin/ory: Makefile
	curl https://raw.githubusercontent.com/ory/meta/master/install.sh | bash -s -- -b .bin ory v0.1.48
	touch .bin/ory

node_modules: package-lock.json
	npm ci
	touch node_modules

gen: .bin/goimports .bin/mockgen  # generates mocks
	./generate-mocks.sh

.DEFAULT_GOAL := help
