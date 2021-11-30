FROM alpine:3.15
COPY znet /bin/znet
RUN mkdir /var/znet
RUN apk --update add ca-certificates
RUN apk add tzdata
ENTRYPOINT ["/bin/znet"]
