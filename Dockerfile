# FROM golang:1.11-alpine3.7
# ADD . /go/src/github.com/xaque208/znet
# RUN go install github.com/xaque208/znet
# CMD znet listen -l :9904 -v
#
FROM alpine:3.11

# Add the binary
COPY bin/linux/znet /bin/znet

ENTRYPOINT ["/bin/znet"]
