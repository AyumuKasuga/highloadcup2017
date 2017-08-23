FROM golang:1.9
WORKDIR /root
ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/go
RUN go get -u github.com/valyala/fasthttp && go get -u github.com/valyala/tcplisten
ADD *.go go/src/srv/
RUN go build srv && go install srv
EXPOSE 80
ENTRYPOINT ./go/bin/srv