FROM golang:onbuild

ENTRYPOINT /go/bin/hydra host

EXPOSE 4444