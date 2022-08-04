# First build application
FROM golang:1.18-alpine

RUN mkdir /app
WORKDIR /app

# Pre-download to enable dependency caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and start build
COPY . .
RUN go build -v -o /app ./...

# Throw away last step and put binary in basic alpine image
FROM alpine:latest

RUN apk add --no-cache iputils setpriv dumb-init && rm -rf /root/.cache

RUN mkdir /app
COPY --from=0 /app/jojo-discord-bot /app/jojo-discord-bot
COPY ./scripts/entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod 755 /app/jojo-discord-bot
RUN chmod 755 /usr/bin/entrypoint.sh

VOLUME ["/app/data", "/app/log"]

# General image informations
LABEL author="Lazy Bytez"
LABEL maintainer="contact@lazybytez.de"

# Open Container annotations
LABEL org.opencontainers.image.title="JOJO Discord Bot"
LABEL org.opencontainers.image.description="An advanced multi-purpose discord bot"
LABEL org.opencontainers.image.vendor ="Lazy Bytez"
LABEL org.opencontainers.image.source="https://github.com/lazybytez/jojo-discord-bot"
LABEL org.opencontainers.image.licenses="AGPL-3.0"

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/usr/bin/entrypoint.sh"]
CMD ["/app/jojo-discord-bot"]
