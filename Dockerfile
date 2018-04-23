# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/downlink
RUN apk update && apk add git
RUN cd $GOPATH/src/github.com/aiotrc/downlink/ && go get -v && go build -v -o /downlink

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /downlink /app/
ENTRYPOINT ["./downlink"]
