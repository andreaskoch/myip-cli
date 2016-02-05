FROM golang:latest
MAINTAINER Andreas Koch <andy@ak7.io>

# Add sources
ADD . /go/src/github.com/andreaskoch/myip
WORKDIR /go/src/github.com/andreaskoch/myip

# Run unit tests with code coverage
RUN make coverage

# Build
RUN make crosscompile
