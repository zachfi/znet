FROM alpine:3.14
COPY bin/linux/znet /bin/znet
ENTRYPOINT ["/bin/znet"]
