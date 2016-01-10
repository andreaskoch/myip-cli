# myip

A small tool that returns your current local or remote IP addresses

**myip** is a tiny, cross-platform command-line utility that just **returns your current IP address** (local or remote) for Linux, Windows, Mac, FreeBSD, NetBSD and OpenBSD.

## Motivation

Of course there are many ways to determine your current local or remote IP address. But all of them (afaik) require some scripting and string extraction and would not work the same way on all platforms.

I created this really simple tool to do just that: **Get the current IP address - on all operating systems**.

## Usage

```bash
myip <action> [options]
```

**Actions**:

- `local`: Get your local IP address
- `remote`: Get your remote IP address

**Options**:

- `-4`: Use IPv4 instead of IPv6 (optional)
- `-select`: Select one or more IPs (optional)
  - `all`: Returns all IP addresses
  - `first`: Returns only the first IP address
  - `last`: Returns only the last IP address
  - `2`: Return only the second IP address
  - `1,2,3`: Return only the first three IP addresses
  - `3,2,1`: Return only the first three IP addresses in reverse order
  - `3`: Return only the third IP address

### Get Help

Get help and usage information:

```bash
myip
```

or

```bash
myip -h
```

### Get the current local IP(s)

Get the current local IP addresses:

```bash
myip local
```

Get only the first local IP address:

```bash
myip local -select first
```

or

```bash
myip local -select 1
```

### Get the current remote IP(s)

Get the current remote IP address:

```bash
myip remote
```

### IPv6 vs. IPv4

myip will only return **IPv6** addresses **by default**. If you want myip to return an IPv4 address you must add the `-4` flag.

Examples:

```bash
myip local -4
myip remote -6
```

## Installation

If you have [go](https://golang.org/) installed:

```bash
git clone git@github.com:andreaskoch/myip.git && cd myip
go run make.go -install
```

or

```bash
go get github.com/andreaskoch/myip
```

## Build

You can use to included `make.go`-script to cross-compile _myip_ for darwin, freebsd, linux, netbsd, openbsd and windows (amd64, arm, arm5, arm6, arm7):

```bash
go run make.go -crosscompile
```

If you don't have golang installed but have [docker](https://www.docker.com/) instead you can use docker to build/cross-compile _myip_:

```bash
git clone git@github.com:andreaskoch/myip.git && cd myip
docker run \
       --rm \
       -v `pwd`:/go/src/github.com/andreaskoch/myip \
       --workdir=/go/src/github.com/andreaskoch/myip \
       golang:latest \
       make crosscompile
```

Or you can extract the binaries from the automatically built [andreaskoch/myip](https://hub.docker.com/r/andreaskoch/allmark/) docker image:

```bash
docker run --rm -v `pwd`:/exchange andreaskoch/myip:latest bash -c "cp -a /go/src/github.com/andreaskoch/myip/bin/* /exchange"
```

↑ This command will copy the binaries from the docker container to your current directory.
