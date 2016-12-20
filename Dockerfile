FROM golang:1.6

# create directory for our code
RUN mkdir -p /go/src/app

# set the default working directory
WORKDIR /go/src/app

# add source to the image
COPY . /go/src/app

# pull dependencies
RUN go-wrapper download

# build and install binary
RUN go-wrapper install

# instruct docker that port 8080 is used
EXPOSE 8080
