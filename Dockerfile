FROM alpine:latest

RUN apk update && apk add ca-certificates
COPY backoff-linux-amd64 /bin/backoff

ENTRYPOINT [ "/bin/backoff" ]
CMD ["--help"]
