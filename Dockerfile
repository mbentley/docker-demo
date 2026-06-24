# build stage: resolve deps and compile
FROM golang:1.26-alpine AS build

WORKDIR /build

# copy source and module definition; go mod tidy resolves and pins deps at build time
COPY go.mod *.go ./
RUN go mod tidy && go build -o /go/bin/app .


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