FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates sqlite

COPY clonis /usr/bin/clonis

COPY static ./static
COPY templates ./templates

RUN chmod +x /usr/bin/clonis

EXPOSE 8080
VOLUME ["/config"]

CMD ["clonis"]