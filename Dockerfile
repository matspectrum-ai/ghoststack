# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=0.0.0
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 go build -ldflags "\
  -X github.com/ghoststack/ghoststack/internal/cli.Version=${VERSION} \
  -X github.com/ghoststack/ghoststack/internal/cli.Commit=${COMMIT} \
  -X github.com/ghoststack/ghoststack/internal/cli.BuildTime=${BUILD_TIME} \
  -s -w" \
  -o /ghost ./cmd/ghost

FROM alpine:3.20

RUN apk add --no-cache \
  ca-certificates \
  iptables \
  nftables \
  wireguard-tools \
  iproute2 \
  bash

COPY --from=builder /ghost /usr/local/bin/ghost

RUN mkdir -p /etc/ghoststack /var/lib/ghoststack && \
  addgroup -S ghoststack && \
  adduser -S -G ghoststack -h /var/lib/ghoststack ghoststack && \
  chown -R ghoststack:ghoststack /etc/ghoststack /var/lib/ghoststack

VOLUME ["/etc/ghoststack", "/var/lib/ghoststack"]

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT ["/usr/local/bin/ghost"]
CMD ["start"]
