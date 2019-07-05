# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.11 as builder

RUN mkdir -p "$GOPATH/src/github.com/I1820/tm"
WORKDIR $GOPATH/src/github.com/I1820/tm
ENV GO111MODULE=on

COPY . .
RUN go build -o /bin/app

FROM alpine:3.9

WORKDIR /bin/

COPY --from=builder /bin/app .

# Comment out to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 1995

# Metadata
ARG BUILD_DATE
ARG BUILD_COMMIT
ARG BUILD_COMMIT_MSG
LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"
LABEL org.i1820.build-date=$BUILD_DATE
LABEL org.i1820.build-commit-sha=$BUILD_COMMIT
LABEL org.i1820.build-commit-msg=$BUILD_COMMIT_MSG

ENTRYPOINT ["/bin/app"]
