FROM docker:dind as docker
FROM golang:1.17.11-buster
# FROM --platform=linux/amd64 golang:1.17.11-buster
RUN apt-get update && apt-get install -y git gcc bash make curl docker-compose 
COPY --from=docker /usr/local/bin/docker /usr/local/bin/docker 
RUN curl -sL https://deb.nodesource.com/setup_16.x | bash && apt-get install -y nodejs
# COPY /Users/macbookpro/dev/shadow/hydra_ory /Users/macbookpro/dev/shadow/hydra
# WORKDIR /Users/macbookpro/dev/shadow/hydra
# RUN apt-get update && apt-get install -y pass gnupg2
WORKDIR /usr
RUN git clone https://github.com/go-delve/delve
WORKDIR /usr/delve
RUN go install github.com/go-delve/delve/cmd/dlv

# ENTRYPOINT [ "gpg2", "–", "gen-key", "&&", "pass", "init", '\$gpg_id' ]

