FROM golang:1.16
RUN mkdir /go/src/web
WORKDIR /go/src/web
ADD . /go/src/web
