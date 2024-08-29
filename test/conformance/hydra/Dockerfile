FROM golang:1.22 AS builder

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
    git gcc bash ssl-cert ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/ory/hydra
RUN mkdir -p ./internal/httpclient

COPY go.mod go.sum ./
COPY internal/httpclient/go.* ./internal/httpclient/

ENV GO111MODULE on
ENV CGO_ENABLED 1

RUN go mod download

COPY . .

RUN go build -tags sqlite -o /usr/bin/hydra

VOLUME /var/lib/sqlite

# Exposing the ory home directory
VOLUME /home/ory

# Declare the standard ports used by hydra (4444 for public service endpoint, 4445 for admin service endpoint)
EXPOSE 4444 4445

RUN mv test/conformance/ssl/ory-ca.* /etc/ssl/certs/ && \
    mv test/conformance/ssl/ory-conformity.crt /etc/ssl/certs/ && \
    mv test/conformance/ssl/ory-conformity.key /etc/ssl/private/ && \
    update-ca-certificates

ENTRYPOINT ["hydra"]
CMD ["serve"]
