FROM golang:1.16
RUN mkdir /go/src/web && go get github.com/pilu/fresh
WORKDIR /go/src/web
ADD . /go/src/web
EXPOSE 8080
CMD ["fresh"]
