FROM golang:1.20 AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/leoparente/opentelemetry-infinity
COPY . .
RUN apt update -y && apt install make build-essential dpkg-dev curl ca-certificates libasan6 wget zip -y
RUN if [[ ! -f "runner/otelcol-contrib" ]] ; then make getotelcol ; fi
RUN CGO_ENABLED=1 make build && mv build/otlpinf /exe

FROM debian:bullseye-slim

RUN mkdir /opt/otlpinf/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
