FROM alpine:3.11
COPY bin/linux/znet /bin/znet
ENTRYPOINT ["/bin/znet"]
