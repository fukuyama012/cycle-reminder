FROM golang:1.11
 
ENV TZ='Asia/Tokyo'
# setup revel
RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel
# setup dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/fukuyama012/cycle-reminder/service/web

EXPOSE 9000
CMD ["revel run -a github.com/fukuyama012/cycle-reminder/service/web prod"]