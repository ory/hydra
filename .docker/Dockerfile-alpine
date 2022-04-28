FROM alpine:3.15

RUN addgroup -S ory; \
    adduser -S ory -G ory -D -H -s /bin/nologin
RUN apk --no-cache --upgrade --latest add ca-certificates

COPY hydra /usr/bin/hydra

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

USER ory

ENTRYPOINT ["hydra"]
CMD ["serve", "all"]
