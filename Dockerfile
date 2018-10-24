FROM golang:1.11-alpine

ARG git_tag
ARG git_commit

RUN apk add --no-cache git build-base curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/ory/hydra

ENV GO111MODULE=on

ADD ./go.mod ./go.mod
ADD ./go.sum ./go.sum

RUN go mod download

ADD . .

RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -X github.com/ory/hydra/cmd.Version=$git_tag -X github.com/ory/hydra/cmd.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/cmd.GitHash=$git_commit" -a -installsuffix cgo -o hydra

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/github.com/ory/hydra/hydra /usr/bin/hydra

ENTRYPOINT ["hydra"]

CMD ["serve", "all"]
