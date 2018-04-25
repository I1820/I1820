# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/dm
RUN apk update && apk add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN cd $GOPATH/src/github.com/aiotrc/dm/ && dep ensure && go build -v -o /dm

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /dm /app/
ENTRYPOINT ["./dm"]
