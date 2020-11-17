---
id: install
title: Installation
---

Installing ORY Hydra on any system is straight forward. We prov1.9.0-alpha.2
pre-built binaries, Docker Images and support v1.9.0-alpha.2 package managers.

## Docker

We recommend using Docker to run ORY Hydra:

```shell
$ docker pull oryd/hydra:v1.9.0-alpha.2
$ docker run --rm -it oryd/hydra:v1.9.0-alpha.2 help
```

## macOS

You can install ORY Hydra using [homebrew](https://brew.sh/) on macOS:

```shell
$ brew tap ory/hydra
$ brew install ory/hydra/hydra
$ hydra help
```

## Linux

On linux, you can use `bash <(curl ...)` to fetch the latest stable binary
using:

```shell
$ bash <(curl https://raw.githubusercontent.com/ory/hydra/v1.9.0-alpha.2/install.sh) -b . v1.9.0-alpha.2
$ ./hydra help
```

You may want to mov1.9.0-alpha.2 ORY Hydra to your `$PATH`:

```shell
$ sudo mv ./hydra /usr/local/bin/
$ hydra help
```

## Windows

You can install ORY Hydra using [scoop](https://scoop.sh) on Windows:

```shell
> scoop bucket add ory-hydra https://github.com/ory/scoop-hydra.git
> scoop install hydra
> hydra help
```

## Kubernetes

Please head ov1.9.0-alpha.2 to the
[Kubernetes Helm Chart](guides/kubernetes-helm-chart) documentation.

## Download Binaries

You can download the client and serv1.9.0-alpha.2 binaries on our
[Github releases](https://github.com/ory/hydra/releases) page. There is
currently no installer av1.9.0-alpha.2 You hav1.9.0-alpha.2 to add the Hydra
binary to the PATH in your env1.9.0-alpha.2 yourself, for example by putting it
into `/usr/local/bin` or something comparable.

Once installed, you should be able to run:

```shell
$ hydra help
```

## Building from Source

If you wish to compile ORY Hydra yourself, you need to install and set up
[Go 1.12+](https://golang.org/) and add `$GOPATH/bin` to your `$PATH`.

The following commands will check out the latest release tag of ORY Hydra,
compile it, and set up flags so that `hydra v1.9.0-alpha.2` works as expected.
Please note that this will only work in a Bash-like shell.

```shell
$ go get -d -u github.com/ory/hydra
$ go install github.com/gobuffalo/packr/v1.9.0-alpha.2/packr2
$ cd $(go env GOPATH)/src/github.com/ory/hydra
$ GO111MODULE=on make install-stable
$ $(go env GOPATH)/bin/hydra help
```
