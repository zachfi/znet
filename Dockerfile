FROM alpine:3.14 as certs
COPY ./bin/linux/znet /bin/znet
RUN chmod 0700 /bin/znet
RUN mkdir /var/znet
RUN apk --update add ca-certificates
RUN apk add libc6-compat
ENTRYPOINT ["/bin/znet"]
