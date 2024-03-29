# start with the official golang 1.17 image based on alpine
FROM golang:1.17-alpine AS build

# create directory for our code & set as the working directory
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# add source to the image
COPY go.mod go.sum *.go /go/src/app/
COPY vendor /go/src/app/vendor/
COPY static /go/src/app/static/
COPY templates /go/src/app/templates/

# pull dependencies and then build the app binary
RUN go get -v -d &&\
  go install -v


# now that our binary is built, let's start from a clean alpine image
# rebased/repackaged base image that only updates existing packages
FROM mbentley/alpine:latest
LABEL maintainer="Matt Bentley <mbentley@mbentley.net>"

# copy binary from build image
COPY --from=build /go/bin/app /app/app

# copy static files, templates, and entrypoint script from build context
COPY static /app/static
COPY templates /app/templates
COPY entrypoint.sh /entrypoint.sh

# set workding directory, entrypoint and command
WORKDIR /app
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/app"]
