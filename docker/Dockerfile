FROM golang:1.20-alpine AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/leoparente/opentelemetry-infinity
COPY . .
RUN apk update \
    && apk add make build-base
RUN make getotelcol \
    && CGO_ENABLED=1 make binary \
    && mv build/otel-infinity /exe

FROM alpine:latest

RUN mkdir /opt/otel-infinity/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
