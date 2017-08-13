# start with the official golang 1.8 image based on alpine
FROM golang:1.8-alpine as build

# create directory for our code & set as the working directory
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# add source to the image
COPY *.go /go/src/app/
COPY vendor /go/src/app/vendor/
COPY static /go/src/app/static/
COPY templates /go/src/app/templates/

# pull dependencies and then build the app binary
RUN go-wrapper download
RUN go-wrapper install


# now that our binary is built, let's start from a clean alpine image
FROM alpine:latest

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
