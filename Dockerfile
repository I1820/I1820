# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.12-alpine3.9 as builder

WORKDIR /app

RUN apk --no-cache add git ca-certificates git gcc g++ libc-dev

COPY . .
RUN go build -o /bin/app

FROM alpine:3.9

WORKDIR /bin/

COPY --from=builder /bin/app .

# Comment out to run the binary in "production" mode:
ENV i1820_link_services_http_debug=false

EXPOSE 1372

# Metadata
ARG BUILD_DATE
ARG BUILD_COMMIT
ARG BUILD_COMMIT_MSG
LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"
LABEL org.i1820.build-date=$BUILD_DATE
LABEL org.i1820.build-commit-sha=$BUILD_COMMIT
LABEL org.i1820.build-commit-msg=$BUILD_COMMIT_MSG

# Comment out to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD ["/bin/app"]
