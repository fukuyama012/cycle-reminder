FROM golang:1.11.2
 
# setup revel
RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel
# setup dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/web

WORKDIR /go/src/web
RUN dep ensure

