# To compile this image manually run:
#
# $ GO111MODULE=on GOOS=linux GOARCH=amd64 go build && docker build -t oryd/hydra:v1.0.0-rc.7_oryOS.10 . && rm hydra
#FROM golang:1.12.7-alpine3.10
FROM golang:1.11

#RUN apk add -U --no-cache ca-certificates

WORKDIR /hydra/
RUN chown 1000 /hydra/

#ADD scripts/docker/install_service.sh install_service.sh

ADD metrics/  metrics/
ADD cmd/  cmd/
ADD tracing/  tracing/
ADD cypress/  cypress/
ADD test/  test/
ADD driver/  driver/
ADD health/  health/
ADD oauth2/  oauth2/
ADD internal/  internal/
#ADD docs/  docs/
ADD jwk/  jwk/
ADD consent/  consent/
ADD sdk/  sdk/
#ADD scripts/  scripts/
ADD client/  client/
ADD x/  x/
ADD vendor/  vendor/
ADD main.go main.go
ADD go.mod go.mod
ADD go.sum go.sum

#RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build

#COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY ./hydra /usr/bin/hydra

USER 1000

#ENTRYPOINT ["hydra"]
#CMD ["serve", "all"]
#CMD ["find", "/go/pkg"]