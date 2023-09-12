## _BUILDING collector binary_ ##

```sh
./builder --config=/mnt/c/tmp/otel_build_config.yaml 
```
(here builder is the otel builder binary , refer : https://github.com/open-telemetry/opentelemetry-collector-builder  )

contents of otel_build_config.yaml (build specs) :
```sh
dist:
  name: otelcol-custom-2
  description: Local OpenTelemetry Collector binary
  output_path: /mnt/c/tmp/otelxxx
exporters:
  - gomod: go.opentelemetry.io/collector/exporter/loggingexporter v0.84.0
  - gomod: go.opentelemetry.io/collector/exporter/otlpexporter v0.84.0
  - gomod: go.opentelemetry.io/collector/exporter/otlphttpexporter v0.84.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.84.0
receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.84.0
extensions:
  - gomod: github.com/shadow0wolf/infa_auth 1.0.5
  - gomod: go.opentelemetry.io/collector/extension/zpagesextension v0.84.0
```

## _RUNNING collector binary_ ##
```sh
/mnt/c/tmp/otelxxx/otelcol-custom-2 --config=/mnt/c/tmp/otel_config.yaml
```

contents of otel_config.yaml (run config ) :

```sh
extensions:
  infa_auth:
    #validation_url: https://pod.ics.dev:444/session-service/api/v1/session/Agent
    validation_url: http://172.20.64.1:9898/session-service/api/v1/session/Agent
    header_key: IDS-AGENT-SESSION-ID
  
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:55681"
        auth:
          authenticator: infa_auth
      http:
        endpoint: "0.0.0.0:55680"
        auth:
          authenticator: infa_auth

exporters:
  file:
    path: "/mnt/c/tmp/otel_logs.txt"
  logging:
    verbosity: detailed  
service:
  extensions: [infa_auth]
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [file,logging]
      
  telemetry:
    logs:
      level: "debug"
```
