GLIDE = $(GOPATH)/bin/glide

SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)

# These are paths used in the docker image
SRCROOT_D = /go/src/github.com/ory-am/hydra
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist

default: build

build: dep
	CGO_ENABLED=0 go build -o $(BUILD_ROOT)/sand

dep:
	go get -u github.com/Masterminds/glide
	if [ ! -d vendor ]; then $(GLIDE) install; fi

dist:
	docker pull golang:1.7.1
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
	           -e UID=`id -u` \
	           -e GID=`id -g` \
	           golang:1.7.1 \
	           make distbuild

distbuild: clean build
	-chown -R $(UID):$(GID) $(SRCROOT)

distdep:
	docker pull golang:1.7.1
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
	           -e UID=`id -u` \
	           -e GID=`id -g` \
	           golang:1.7.1 \
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

.PHONY: bin default dep updatedep dist distdep clean
