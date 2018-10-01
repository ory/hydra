GLIDE = $(GOPATH)/bin/glide

SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)

# These are paths used in the docker image
SRCROOT_D = /go/src/github.com/ory/hydra
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist

REVISION = $$(git rev-parse --short HEAD)
VERSION = $$(git name-rev --tags --name-only $(REVISION))

LD_FLAGS ?= -ldflags "-X github.com/ory/hydra/cmd.GitHash=$(REVISION) -X github.com/ory/hydra/cmd.Version=$(VERSION)"

default: build

build: dep
	CGO_ENABLED=0 go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand

build-osx: dep
	CGO_ENABLED=0 GOOS=darwin go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand-osx
	tar cvzf $(BUILD_ROOT)/sand-$(VERSION)-osx.tgz $(BUILD_ROOT)/sand-osx

build-linux: dep
	CGO_ENABLED=0 GOOS=linux go build $(LD_FLAGS) \
		-o $(BUILD_ROOT)/sand-linux
	tar cvzf $(BUILD_ROOT)/sand-$(VERSION)-linux.tgz $(BUILD_ROOT)/sand-linux

dep:
	go get -u github.com/Masterminds/glide
	if [ ! -d vendor ]; then $(GLIDE) install; fi

dist:
	docker pull golang:1.11.0
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
	           -e UID=`id -u` \
	           -e GID=`id -g` \
	           golang:1.11.0 \
	           make distbuild

distbuild: clean build build-osx build-linux
	-chown -R $(UID):$(GID) $(SRCROOT)

distdep:
	docker pull golang:1.11.0
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
	           -e UID=`id -u` \
	           -e GID=`id -g` \
	           golang:1.11.0 \
	           make updatedep

updatedep:
	if [ -f glide.yaml ]; then mv glide.yaml .glide.yaml.backup; fi
	if [ -f glide.lock ]; then mv glide.lock .glide.lock.backup; fi
	go get -u github.com/Masterminds/glide
	${GLIDE} create
	${GLIDE} up

clean:
	if [ -d $(BUILD_ROOT_D) ]; then rm -rf $(BUILD_ROOT_D); fi
	-chown -R $(UID):$(GID) $(SRCROOT)
	if [ -d $(SRCROOT)/vendor ]; then rm -rf $(SRCROOT)/vendor; fi

.PHONY: bin default build build-osx build-linux dep updatedep dist distbuild distdep clean
