# Start from the latest golang base image
FROM golang:alpine AS builder

RUN apk --no-cache add git ca-certificates git gcc g++ libc-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /I1820

FROM alpine:latest

# Metadata
ARG BUILD_DATE
ARG BUILD_COMMIT
ARG BUILD_COMMIT_MSG
LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"
LABEL org.i1820.build-date=$BUILD_DATE
LABEL org.i1820.build-commit-sha=$BUILD_COMMIT
LABEL org.i1820.build-commit-msg=$BUILD_COMMIT_MSG

WORKDIR /root/

COPY --from=builder /I1820 .

ENTRYPOINT ["./I1820"]
