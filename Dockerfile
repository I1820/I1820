# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/uplink
RUN apk update && apk add git
RUN cd $GOPATH/src/github.com/aiotrc/uplink/ && go get -v && go build -v -o /uplink

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /uplink /app/
ENTRYPOINT ["./uplink"]
