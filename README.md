# myip

A small tool that returns your current local IP address

**myip** is a tiny, cross-platform command-line utility that just **returns your current IP address** (local or remote) for Linux, Windows, Mac, FreeBSD, NetBSD and OpenBSD.

## Motivation

Of course there are many ways to determine your current local or remote IP address. But all of them (afaik) require some scripting and string extraction and would not work on multiple platforms.

I created this really simple tool to do just **one task - on all operating systems alike**.

## Usage

Get the current IPv6 address:

```bash
myip
```

Get the current IPv4 address:

```bash
myip -4
```

myip will only return IPv6 addresses by default. If you want myip to return an IPv4 address you must add the `-4` flag.

## Installation

If you have [go](https://golang.org/) installed:

```bash
go get github.com/andreaskoch/myip
```

## Build

You can use [gox](https://github.com/mitchellh/gox) to compile _myip_ for darwin, freebsd, linux, netbsd, openbsd and windows (x86, amd64, arm):

```bash
make install
```

or

```bash
go get github.com/mitchellh/gox
gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
```

If you don't have golang installed but have [docker](https://www.docker.com/) you can use docker to build/cross-compile _myip_:

```bash
git clone git@github.com:andreaskoch/myip.git && cd myip
docker run \
       --rm \
       -v `pwd`:/go/src/github.com/andreaskoch/myip \
       --workdir=/go/src/github.com/andreaskoch/myip \
       golang:latest \
       make install
```

Or you can extract the binaries from the automatically built [andreaskoch/myip](https://hub.docker.com/r/andreaskoch/allmark/) docker image:

```bash
docker run --rm -v `pwd`:/exchange andreaskoch/myip:latest bash -c "cp -a /go/src/github.com/andreaskoch/myip/bin/* /exchange"
```

↑ This command will copy the binaries from the docker container to your current directory.
