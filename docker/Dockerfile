FROM golang:1.20-alpine AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/leoparente/opentelemetry-infinity
COPY . .
RUN apk update && apk add make build-base dpkg-dev curl
RUN if [ ! -f "runner/otelcol-contrib" ] ; then make getotelcol ; fi
RUN make build && mv build/otlpinf /exe

FROM alpine:3.18.5

RUN mkdir /opt/otlpinf/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /

ENTRYPOINT ["/exe"]
