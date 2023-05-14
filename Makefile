VERSION ?= $(shell runner/otelcol-contrib -v | cut -d ' ' -f 3)
BUILD_DIR = build
CGO_ENABLED ?= 0
DOCKERHUB_REPO = ghcr.io/leoparente
GOARCH ?= $(shell dpkg-architecture -q DEB_BUILD_ARCH)
COMMIT_HASH = $(shell git rev-parse --short HEAD)
INF_LATEST_RELEASE=$(shell curl -L -s -H 'Accept: application/json' https://github.com/leoparente/opentelemetry-infinity/releases/latest)
INF_LATEST_VERSION=$(shell echo ${INF_LATEST_RELEASE} | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')

getotelcol:
	wget -O /tmp/otelcol-contrib-$(GOARCH)$(GOARM).zip https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/$(INF_LATEST_VERSION)/otelcol-contrib-$(GOARCH)$(GOARM).zip
	unzip /tmp/otelcol-contrib-$(GOARCH)$(GOARM).zip -d /tmp/
	mv /tmp/otelcol-contrib runner/otelcol-contrib
	rm -rf /tmp/otelcol-contrib*

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o ${BUILD_DIR}/otlpinf cmd/main.go
  
container:
	docker build --no-cache \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:develop \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:develop-$(COMMIT_HASH) \
	  -f docker/Dockerfile .

release:
	docker build --no-cache \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:latest \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:$(VERSION) \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:$(VERSION)-$(COMMIT_HASH) \
	  -f docker/Dockerfile .
