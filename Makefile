GLIDE = $(GOPATH)/bin/glide
GO_VERSION = 1.13.5

SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)

# These are paths used in the docker image
SRCROOT_D = /go/src/github.com/ory/hydra
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist

REVISION = $$(git rev-parse --short HEAD)
VERSION = $$(git name-rev --tags --name-only $(REVISION))

LD_FLAGS ?= -ldflags "-X github.com/ory/hydra/cmd.GitHash=$(REVISION) -X github.com/ory/hydra/cmd.Version=$(VERSION)"

default: build

build:
	CGO_ENABLED=0 go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand

build-osx:
	CGO_ENABLED=0 GOOS=darwin go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand-osx
	tar cvzf $(BUILD_ROOT)/sand-$(VERSION)-osx.tgz $(BUILD_ROOT)/sand-osx

build-linux:
	CGO_ENABLED=0 GOOS=linux go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand-linux
	tar cvzf $(BUILD_ROOT)/sand-$(VERSION)-linux.tgz $(BUILD_ROOT)/sand-linux

dist:
	docker pull golang:$(GO_VERSION)
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
	           -e UID=`id -u` \
	           -e GID=`id -g` \
	           golang:$(GO_VERSION) \
	           make distbuild

distbuild: clean build build-osx build-linux
	-chown -R $(UID):$(GID) $(SRCROOT)

clean:
	if [ -d $(BUILD_ROOT_D) ]; then rm -rf $(BUILD_ROOT_D); fi
	-chown -R $(UID):$(GID) $(SRCROOT)
	if [ -d $(SRCROOT)/vendor ]; then rm -rf $(SRCROOT)/vendor; fi

.PHONY: bin default build build-osx build-linux dist distbuild clean
