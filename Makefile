VERSION = $(shell cat cmd/VERSION)
BUILD_DIR = build
CGO_ENABLED ?= 0
DOCKERHUB_REPO = ghcr.io/leoparente
GOARCH ?= $(shell dpkg-architecture -q DEB_BUILD_ARCH)
COMMIT_HASH = $(shell git rev-parse --short HEAD)

define compile_service
  echo "VERSION: $(VERSION)"
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=mod -ldflags "-extldflags "-static" -X 'github.com/leoparente/opentelemetry-infinity/buildinfo.version=$(VERSION)'" -o ${BUILD_DIR}/otel-infinity cmd/main.go
endef

getotelcol:
	wget -o /tmp/otelcol-contrib.tar.gz https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.75.0/otelcol-contrib_0.75.0_linux_amd64.tar.gz
	tar -xvzf /tmp/otelcol-contrib.tar.gz -C /tmp/
	mv /tmp/otelcol-contrib runner/otelcol-contrib
	rm -rf /tmp/otelcol-contrib*

binary:
	$(call compile_service)
  
container:
	docker build --no-cache \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:$(REF_TAG) \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:$(VERSION) \
	  --tag=$(DOCKERHUB_REPO)/opentelemetry-infinity:$(VERSION)-$(COMMIT_HASH) \
	  -f docker/Dockerfile .
