FROM alpine:3.13
COPY bin/linux/znet /bin/znet
ENTRYPOINT ["/bin/znet"]
