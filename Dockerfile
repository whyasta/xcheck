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
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o job ./cmd/job
# ================

# ================ Start running app
FROM alpine:3.20
# Install Dependencies
RUN apk --no-cache add libaio libnsl libc6-compat curl supervisor

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/job .

# COPY certs/*.crt /etc/ssl/certs/
RUN touch .env
RUN mkdir /app/config
RUN mkdir /var/log/supervisor
# COPY config/development.yml config/development.yml
# COPY config/production.yml config/production.yml
COPY app.conf /etc/supervisor/conf.d/
COPY start.sh .
RUN chown -R root:root /app && chmod -R ug+rwx /app

EXPOSE 9052
# CMD ["/app/main", "-e", "production"]
ENTRYPOINT ["/app/start.sh"]
# ================ End running app
