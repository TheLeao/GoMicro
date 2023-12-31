FROM alpine:latest

RUN mkdir /app

COPY mail-service app
COPY templates templates

CMD ["/app/mail-service"]