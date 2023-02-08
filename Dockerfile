FROM alpine:3.17 as certs
COPY ./bin/linux/znet /bin/znet
RUN chmod 0700 /bin/znet
RUN mkdir /var/znet
RUN apk --update add ca-certificates
RUN apk add libc6-compat
RUN apk add tzdata
ENTRYPOINT ["/bin/znet"]
