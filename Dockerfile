FROM golang:1.7

ADD . /go/src/github.com/ory-am/hydra
WORKDIR /go/src/github.com/ory-am/hydra

RUN go get github.com/Masterminds/glide
RUN glide install
RUN go install github.com/ory-am/hydra

ENTRYPOINT /go/bin/hydra host

EXPOSE 4444