FROM golang:1.6
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/trading-post
COPY . /go/src/github.com/octoblu/trading-post

RUN env CGO_ENABLED=0 go build -o trading-post -a -ldflags '-s' .

CMD ["./trading-post"]
