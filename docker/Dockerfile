FROM golang:1.4

ADD . /go/src/github.com/fmonniot/webhook-listener

WORKDIR /go/src/github.com/fmonniot/webhook-listener/dns-listener

RUN go get && go build

ENTRYPOINT ["/go/src/github.com/fmonniot/webhook-listener/dns-listener/dns-listener"]
