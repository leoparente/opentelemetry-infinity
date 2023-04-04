# opentelemetry-infinity

Opentelemetry Infinity provison [otel-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) instances through a simple REST API using a policy mechanism. Each policy spin up a new `otelcol-contrib` process running the configuration provided by the policy.

## Project premises
**1. Single binary**: `otlpinf` embeds `otelcol-contrib` in its binary. Therefore, only one static binary is provided.

**2. No persistence**: `opentelemetry-infinity` stores data in memory and in temporary files only. This adds a new paradigm to `opentelemetry-collector` that is expected to run over a persisted config file as default. If you are looking for a opentelemetry orchestrator as the way it was planned to perform, you should try the official [opentelemetry-operator](https://github.com/open-telemetry/opentelemetry-operator).

**3. Compatibility**: `opentelemetry-infinity` is basically a wrapper over the official `opentelemetry-collector` which has not released a version `1.0` yet, i.e., breaking changes are expected. Any changes that occurs on its CLI will be reflected in this project.

## Docker Image
You can download and run `otlpinf` using docker image:
```
docker run --net=host ghcr.io/leoparente/opentelemetry-infinity run
```

## REST API
The default `otlpinf` address is `localhost:10222`. to change that you can specify host and port when starting `otlpinf`:
```sh
./otlpinf run -a {host} -p {port}
```

### Routes (v1)
`otlpinf` is aimed to be simple and straightforward. 

#### Get runtime and capabilities information

<details>
 <summary><code>GET</code> <code><b>/api/v1/status</b></code> <code>(gets otlpinf runtime data)</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10222/api/v1/status
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/api/v1/capabilities</b></code> <code>(gets otelcol-contrib capabilities)</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10222/api/v1/capabilities
> ```

</details>

#### Policies Management

<details>
 <summary><code>GET</code> <code><b>/api/v1/policies</b></code> <code>(gets all existing policies)</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | JSON array containing all applied policy names                      |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10222/api/v1/policies
> ```

</details>


<details>
 <summary><code>POST</code> <code><b>/api/v1/policies</b></code> <code>(Creates a new policy)</code></summary>

##### Parameters

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | None      |  required | YAML object             | yaml format specified in [Policy RFC](#policy-rfc-v1) 


##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/x-yaml" --data @post.yaml http://localhost:10222/api/v1/policies
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/api/v1/policies/{policy_name}</b></code> <code>(gets information of a specific policy)</code></summary>

##### Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> |   `policy_name`   |  required | string         | The unique policy name              |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10222/api/v1/policies/my_policy
> ```

</details>

<details>
 <summary><code>DELETE</code> <code><b>/api/v1/policies/{policy_name}</b></code> <code>(delete a existing policy)</code></summary>

##### Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> |   `policy_name`   |  required | string         | The unique policy name              |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X DELTE -H "Content-Type: application/json" http://localhost:10222/api/v1/policies/my_policy
> ```

</details>

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
