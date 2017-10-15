# Build
FROM golang:1.9 as build

WORKDIR /go/src/github.com/royvandewater/trading-post
COPY . /go/src/github.com/royvandewater/trading-post

RUN env CGO_ENABLED=0 go build -o trading-post -a -ldflags '-s' .

# Entrypoint
FROM centurylink/ca-certs

EXPOSE 80

COPY --from=build /go/src/github.com/royvandewater/trading-post/trading-post /trading-post
ADD html html
ENTRYPOINT ["./trading-post"]
