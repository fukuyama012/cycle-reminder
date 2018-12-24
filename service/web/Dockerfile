FROM golang:1.11.2
 
# setup revel
RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel
# setup dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/fukuyama012/cycle-reminder/service/web

WORKDIR /go/src/github.com/fukuyama012/cycle-reminder/service/web
RUN dep ensure

