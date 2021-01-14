BUILDER_IMAGE = 899991151204.dkr.ecr.us-east-1.amazonaws.com/golang:alpine

SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)

# These are paths used in the docker image
SRCROOT_D = /go/src/github.com/ory/hydra
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist

MAJOR_VERSION = 0
MINOR_VERSION = 9
PATCH_VERSION = 16

REVISION ?= $$(git rev-parse --short HEAD)
VERSION ?= $(MAJOR_VERSION).$(MINOR_VERSION).$(PATCH_VERSION)

BUILD_NUMBER ?= x
BUILD_IDENTIFIER = _${BUILD_NUMBER}

LD_FLAGS ?= -ldflags "-X github.com/ory/hydra/cmd.GitHash=$(REVISION) -X github.com/ory/hydra/cmd.Version=$(VERSION)"

build-cli:
	go build $(LD_FLAGS) -o $(BUILD_ROOT)/sand

docker.build.cli: clean
	docker build -f Dockerfile-cli \
		-t sand$(BUILD_IDENTIFIER) \
		--build-arg VERSION=$(VERSION) \
		--build-arg REVISION=$(REVISION) .
	docker create -it --name tocopy-sand$(BUILD_IDENTIFIER) sand$(BUILD_IDENTIFIER) bash
	docker cp tocopy-sand$(BUILD_IDENTIFIER):$(BUILD_ROOT_D)/sand $(SRCROOT)/
	docker rm -f tocopy-sand$(BUILD_IDENTIFIER)
	docker rmi -f sand$(BUILD_IDENTIFIER)

default: build

build:
	echo "building with: `go version`"; \
  	GO111MODULE=on CGO_ENABLED=1 go build $(LD_FLAGS) -o $(BUILD_ROOT)/sand \
	&& echo "build successful. now checking goboring symbols exists..." && go tool nm $(BUILD_ROOT)/sand | grep _Cfunc__goboringcrypto_ > /dev/null

build-osx:
	CGO_ENABLED=0 GOOS=darwin go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand-osx
	tar cvzf $(BUILD_ROOT)/sand-$(VERSION)-osx.tgz $(BUILD_ROOT)/sand-osx

build-linux:
	echo "building with: `go version`"; \
	GO111MODULE=on CGO_ENABLED=1 GOOS=linux go build $(LD_FLAGS) -o $(BUILD_ROOT)/sand-linux \
	&& echo "build successful. now checking goboring symbols exists..." && go tool nm $(BUILD_ROOT)/sand-linux | grep _Cfunc__goboringcrypto_ > /dev/null
	tar cvzf $(BUILD_ROOT)/sand-$(VERSION)-linux.tgz $(BUILD_ROOT)/sand-linux

distbuild: clean build build-osx build-linux

clean:
	rm -f sand
	if [ -d $(BUILD_ROOT_D) ]; then rm -rf $(BUILD_ROOT_D); fi
	if [ -d $(SRCROOT)/vendor ]; then rm -rf $(SRCROOT)/vendor; fi

.PHONY: bin default build build-osx build-linux dist distbuild clean
