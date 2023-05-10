FROM golang:1.20 AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/leoparente/opentelemetry-infinity
COPY . .
RUN apt update -y \
    && apt install make build-essential dpkg-dev curl ca-certificates libasan6 wget gzip -y
RUN make getotelcol \
    && rm -rf build \
    && CGO_ENABLED=1 make build \
    && mv build/otel-infinity /exe

FROM debian:bullseye-slim

RUN mkdir /opt/otel-infinity/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
