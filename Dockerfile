# build stage: resolve deps and compile
FROM golang:1.26-alpine AS build

WORKDIR /build

# download pinned deps first so this layer is cached until go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy source and compile
COPY *.go ./
RUN go build -o /go/bin/app .


# runtime stage: minimal alpine image
FROM mbentley/alpine:latest
LABEL maintainer="Matt Bentley <mbentley@mbentley.net>"

# copy binary from build stage
COPY --from=build /go/bin/app /app/app

# copy static files, templates, and entrypoint script from build context
COPY static /app/static
COPY templates /app/templates
COPY entrypoint.sh /entrypoint.sh

WORKDIR /app
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/app"]