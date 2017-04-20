FROM golang:1.8-alpine

RUN apk add --no-cache git
RUN go get github.com/Masterminds/glide
ADD . /go/src/github.com/ory-am/hydra
WORKDIR /go/src/github.com/ory-am/hydra

ADD ./glide.yaml ./glide.yaml
ADD ./glide.lock ./glide.lock
RUN glide install --skip-test -v

# ADD . .
RUN go install .

ENTRYPOINT /go/bin/hydra host

EXPOSE 4444
