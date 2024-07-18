# ================
FROM golang:1.22.5-alpine3.20 AS builder

LABEL maintainer="Bigmind"

RUN \
    set -ex && \
    apk update && \
    apk add make libc-dev gcc libtool musl-dev ca-certificates dumb-init

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main ./cmd/xcheck
# ================

# ================ Start running app
FROM alpine:3.16
# Install Dependencies
RUN apk --no-cache add libaio libnsl libc6-compat curl

WORKDIR /app
COPY --from=builder /app/main .

# COPY certs/*.crt /etc/ssl/certs/

RUN mkdir /app/config
# COPY config/development.yml config/development.yml
# COPY config/production.yml config/production.yml
COPY start.sh .
EXPOSE 9052
CMD ["/app/main", "-e", "production"]
# ================ End running app
