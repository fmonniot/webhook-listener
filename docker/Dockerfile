FROM golang:1.3
ADD . /go/src/github.com/fmonniot/dns-webhook-listener
WORKDIR /go/src/github.com/fmonniot/dns-webhook-listener/dns-listener
RUN go get && go build
ENTRYPOINT ["/go/src/github.com/fmonniot/dns-webhook-listener/dns-listener/dns-listener"]
