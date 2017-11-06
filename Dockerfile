FROM alpine:3.2

RUN apk add --update ca-certificates # Certificates for SSL

ADD dist/hydra-linux-amd64 /go/bin/hydra

ENTRYPOINT ["/go/bin/hydra", "host"]
