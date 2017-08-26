FROM golang:latest
WORKDIR /root
ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/go
RUN go get -u github.com/valyala/fasthttp && go get -u github.com/valyala/tcplisten && go get -u github.com/hashicorp/golang-lru
RUN go version
ADD *.go go/src/srv/
RUN go build srv && go install srv
EXPOSE 80
ENTRYPOINT ./go/bin/srv