# Pure Storage Flash Array Receiver

| Status                   |                     |
| ------------------------ |---------------------|
| Stability                | [in development]    |
| Supported pipeline types | metrics             |
| Distributions            | [contrib]           |

The Pure Storage Flash Array receiver, receives metrics from Pure Storage internal services hosts.

Supported pipeline types: metrics

## Configuration

The following settings are required:

 -  

The following setting are optional:

 - 
 
Example:

```yaml
receivers:
  purefa:
    confighttp.Client
    arrays:
      address:
      token:
    settings:
      metrics_x_collection_interval:
```

The full list of settings exposed for this receiver are documented [here](./config.go)
with detailed sample configurations [here](./testdata/config.yaml).

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
