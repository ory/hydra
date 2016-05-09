FROM golang:1.6-onbuild

ENTRYPOINT /go/bin/hydra host

EXPOSE 4444