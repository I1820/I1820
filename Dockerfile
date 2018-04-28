MAINTAINER parham.alvani@gmail.com
# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/pm

RUN apk update && apk add git curl
RUN curl -L -s "https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64" -o $GOPATH/bin/dep && chmod +x $GOPATH/bin/dep

RUN cd $GOPATH/src/github.com/aiotrc/pm/ && $GOPATH/bin/dep ensure && go build -v -o /pm

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /pm /app/
ENTRYPOINT ["./pm"]
