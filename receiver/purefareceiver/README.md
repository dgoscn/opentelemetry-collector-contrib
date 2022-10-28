# Pure Storage FlashArray Receiver

| Status                   |                     |
| ------------------------ |---------------------|
| Stability                | [in-development]    |
| Supported pipeline types | metrics             |
| Distributions            | [contrib]           |

The Pure Storage FlashArray receiver, receives metrics from Pure Storage internal services hosts.

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
    endpoint:
    arrays:
      address:
      token:
    settings:
      metrics_array_collection_interval:
      metrics_host_collection_interval:
      metrics_volume_collection_interval:
      metrics_pods_collection_interval:
      metrics_directories_collection_interval:
```

The full list of settings exposed for this receiver are documented [here](./config.go)
with detailed sample configurations [here](./testdata/config.yaml).

[in-development]: https://github.com/open-telemetry/opentelemetry-collector#in-development
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
