.PHONY: gotags deps vendor build

GO15VENDOREXPERIMENT=1

default: build

deps:
	go get -u github.com/jstemmer/gotags
	go get -u github.com/tools/godep

gotags:
	gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .

vendor:
	rm -rf Godeps
	rm -rf vendor
	godep save ./...

build:
	go build -o hydra-host github.com/ory-am/hydra/cli/hydra-host
	go build -o hydra-signin github.com/ory-am/hydra/cli/hydra-signin
	go build -o hydra-signup github.com/ory-am/hydra/cli/hydra-signup
