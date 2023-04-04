# opentelemetry-infinity

Opentelemetry Infinity provison [otel-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) instances through a simple REST API using a policy mechanism. Each policy spin up a new `otelcol-contrib` process running the configuration provided by the policy.

## Project premises
**1. Single binary**: `otlpinf` embeds `otelcol-contrib` in its binary. Therefore, only one static binary is provided.

**2. No persistence**: `opentelemetry-infinity` stores data in memory and in temporary files only. This adds a new paradigm to `opentelemetry-collector` that is expected to run over a persisted config file as default. If you are looking for a opentelemetry orchestrator as the way it was planned to perform, you should try the official [opentelemetry-operator](https://github.com/open-telemetry/opentelemetry-operator).

**3. Compatibility**: `opentelemetry-infinity` is basically a wrapper over the official `opentelemetry-collector` which has not released a version `1.0` yet, i.e., breaking changes are expected. Any changes that occurs on its CLI will be reflected in this project.

## REST API
The default `otlpinf` address is `localhost:10222`. to change that you can specify host and port when starting `otlpinf`:
```sh
./otlinf run -a {host} -p {port}
```

## Docker Image
You can download and run `otlpinf` using docker image:
```
docker run --net=host ghcr.io/leoparente/opentelemetry-infinity run
```

## Policy RFC (v1) 

```yaml
policy_name:
  #Optional
  feature_gates:
    - confmap.expandEnabled
    - exporter.datadog.hostname.preview
  #TODO: set not implemented yet
  set:
    processors.batch.timeout: 2s
  #Required: Same configuration that you would use inside the config file passed to a otel-collector
  config:
    receivers:
      otlp:
        protocols:
          http:
          grpc: 
 
    exporters:
      logging:
        loglevel: debug
      
    service:
      pipelines:
        metrics:
          receivers:
          - otlp
          exporters:
          - logging
```
