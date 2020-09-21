# To compile this image manually run:
#
# $ GO111MODULE=on GOOS=linux GOARCH=amd64 go build -tags sqlite && docker build -t oryd/hydra:v1.0.0-rc.7_oryOS.10 . && rm hydra
FROM alpine:3.11

RUN apk add -U --no-cache ca-certificates

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/nsswitch.conf /etc/nsswitch.conf
COPY hydra /usr/bin/hydra

USER 1000

ENTRYPOINT ["hydra"]
CMD ["serve", "all"]
