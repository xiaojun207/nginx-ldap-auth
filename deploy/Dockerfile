FROM alpine:latest AS production
#作者
MAINTAINER xiaojun "xiaojun207@126.com"

RUN mkdir /app
WORKDIR /app

COPY App /app/App
COPY cfg.json /app/cfg.json
COPY views /app/views

ENTRYPOINT ["sh", "/app/App"]

