# opentelemetry-infinity

Opentelemetry Infinity provison [otel-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) instances through a simple REST API using a policy mechanism.

## Project premises
1. Single binary: `otlpinf` embeds `otelcol-contrib` in its binary. Therefore, only one static binary is required.
2. No persistence: `opentelemetry-infinity` stores data in memory and in temporary files only. This adds a new paradigm to `opentelemetry-collector` that is expected to run over a persisted config file. If you are looking for a opentelemetry orchestrator as the wait it was planned, you should look the official [opentelemetry-operator](https://github.com/open-telemetry/opentelemetry-operator).
3. Compatibility: `opentelemetry-infinity` is basically a wrapper over the official `opentelemetry-collector` which has not released a version 1.0 yet, i.e., breaking changes are expected. Any changes that occurs on its CLI interface will be reflected in this project.

## Policy RFC (v1) 

```

```
