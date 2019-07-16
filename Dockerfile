# To compile this image manually run:
#
# $ GO111MODULE=on GOOS=linux GOARCH=amd64 go build && docker build -t oryd/hydra:v1.0.0-rc.7_oryOS.10 . && rm hydra
#FROM golang:1.12.7-alpine3.10
FROM golang:1.11

#RUN apk add -U --no-cache ca-certificates

WORKDIR /go/src/github.com/callstats-io/hydra/
#ADD scripts/docker/install_service.sh install_service.sh

ADD metrics/  /go/src/github.com/callstats-io/hydra/metrics/
ADD cmd/  /go/src/github.com/callstats-io/hydra/cmd/
ADD tracing/  /go/src/github.com/callstats-io/hydra/tracing/
ADD cypress/  /go/src/github.com/callstats-io/hydra/cypress/
ADD test/  /go/src/github.com/callstats-io/hydra/test/
ADD driver/  /go/src/github.com/callstats-io/hydra/driver/
ADD health/  /go/src/github.com/callstats-io/hydra/health/
ADD oauth2/  /go/src/github.com/callstats-io/hydra/oauth2/
ADD internal/  /go/src/github.com/callstats-io/hydra/internal/
#ADD docs/  /go/src/github.com/callstats-io/hydra/docs/
ADD jwk/  /go/src/github.com/callstats-io/hydra/jwk/
ADD consent/  /go/src/github.com/callstats-io/hydra/consent/
ADD sdk/  /go/src/github.com/callstats-io/hydra/sdk/
#ADD scripts/  /go/src/github.com/callstats-io/hydra/scripts/
ADD client/  /go/src/github.com/callstats-io/hydra/client/
ADD x/  /go/src/github.com/callstats-io/hydra/x/
ADD vendor/  /go/src/github.com/callstats-io/hydra/vendor/
ADD main.go /go/src/github.com/callstats-io/hydra/main.go
ADD go.mod /go/src/github.com/callstats-io/hydra/go.mod
ADD go.sum /go/src/github.com/callstats-io/hydra/go.sum

RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go get
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build
RUN ls ./hydra

#COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY ./hydra /usr/bin/hydra

USER 1000

ENTRYPOINT ["hydra"]
CMD ["serve", "all"]
