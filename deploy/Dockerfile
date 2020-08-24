FROM alpine:latest AS production

RUN mkdir /app
WORKDIR /app

ARG exePath

COPY $exePath /app/App
COPY cfg.json /app/cfg.json
COPY views /app/views

COPY deploy/entrypoint.sh /app/entrypoint.sh
ENTRYPOINT ["sh", "/app/entrypoint.sh"]

