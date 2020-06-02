FROM alpine:3.12
COPY bin/linux/znet /bin/znet
ENTRYPOINT ["/bin/znet"]
