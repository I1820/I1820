# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/dm
RUN apk update && apk add git
RUN cd $GOPATH/src/github.com/aiotrc/dm/ && go get -v && go build -v -o /dm

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /dm /app/
ENTRYPOINT ["./dm"]
