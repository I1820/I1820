# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.11.1 as builder

RUN mkdir -p "$GOPATH/src/github.com/I1820/pm"
WORKDIR $GOPATH/src/github.com/I1820/pm

COPY . .
RUN dep ensure
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# Comment out to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 8080

# Metadata
LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"
LABEL org.i1820.build-date=$BUILD_DATE
LABEL org.i1820.build-commit-sha=$BUILD_COMMIT
LABEL org.i1820.build-commit-msg=$BUILD_COMMIT_MSG

# Comment out to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD ["/bin/app"]
