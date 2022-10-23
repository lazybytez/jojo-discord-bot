# First build application
FROM golang:1.19-alpine

ARG app_version="edge"
ARG build_commit_sha=""

ENV APP_VERSION=$app_version
ENV BUILD_COMMIT_SHA=$build_commit_sha

RUN mkdir -p /app
WORKDIR /app

# Install necessary dependencies
RUN apk update && apk upgrade && apk add g++

# Pre-download to enable dependency caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and start build
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init
RUN go build -ldflags "-X github.com/lazybytez/jojo-discord-bot/build.Version=${APP_VERSION} -X github.com/lazybytez/jojo-discord-bot/build.CommitSHA=${BUILD_COMMIT_SHA}" -v -o /app ./...

# Throw away last step and put binary in basic alpine image
FROM alpine:latest

RUN apk add --no-cache iputils setpriv dumb-init && rm -rf /root/.cache

RUN mkdir -p /app
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
LABEL org.opencontainers.image.vendor="Lazy Bytez"
LABEL org.opencontainers.image.source="https://github.com/lazybytez/jojo-discord-bot"
LABEL org.opencontainers.image.licenses="AGPL-3.0"

# Default port for WebAPI
EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/usr/bin/entrypoint.sh"]
CMD ["/app/jojo-discord-bot"]
