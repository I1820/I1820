# Build stage
FROM golang:alpine AS build-env
COPY . $GOPATH/src/github.com/I1820/link
RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR $GOPATH/src/github.com/I1820/link/
RUN dep ensure && go build -v -o /link

# Final stage
FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build-env /link /app/
ENTRYPOINT ["./link"]
