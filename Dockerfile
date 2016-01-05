FROM golang:latest
MAINTAINER Andreas Koch <andy@ak7.io>

# Add sources
ADD . /go/src/github.com/andreaskoch/myip
WORKDIR /go/src/github.com/andreaskoch/myip

# Build
RUN make crosscompile && make clean
