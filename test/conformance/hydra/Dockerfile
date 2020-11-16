FROM golang:1.15-buster AS builder

RUN apt-get update && \
  apt-get install -y git gcc bash ssl-cert ca-certificates

WORKDIR /go/src/github.com/ory/hydra

ADD go.mod go.mod
ADD go.sum go.sum

ENV GO111MODULE on
ENV CGO_ENABLED 1

RUN go mod download

ADD . .

RUN make pack && go build -tags sqlite -o /usr/bin/hydra

VOLUME /var/lib/sqlite

# Exposing the ory home directory to simplify passing in hydra configuration (e.g. if the file $HOME/.hydra.yaml
# exists, it will be automatically used as the configuration file).
VOLUME /home/ory

# Declare the standard ports used by hydra (4433 for public service endpoint, 4434 for admin service endpoint)
EXPOSE 4444 4445

RUN mv test/conformance/ssl/ory-ca.* /etc/ssl/certs/
RUN mv test/conformance/ssl/ory-conformity.crt /etc/ssl/certs/
RUN mv test/conformance/ssl/ory-conformity.key /etc/ssl/private/
RUN update-ca-certificates

ENTRYPOINT ["hydra"]
CMD ["serve"]
