FROM alpine:latest

RUN mkdir /app

COPY authentication-service app

CMD ["app/authentication-service"]