## 1. PROBLEM STATEMENT  ##
 **Process only valid telemetry requests coming from valid secure-agents.**
 
OTEL Collector (deployed in enterprise cloud) will recieve HTTP/gRPC telemetry-requests (telemetry-signals) from informatica-secure-agent via internet , thus it is required to filter malicious requests (DOS attack) and allow the data from validated telemetry-signals to be exported to downtream component.    
## 2. GUIDELINES PROVIDED BY OTEL ON CUSTOM AUTHENTICATOR IMPLEMENTATION ##
(ref :[official otel doc] https://opentelemetry.io/docs/collector/custom-auth/)

The OpenTelemetry Collector allows receivers and exporters to be connected to authenticators, providing a way to both authenticate incoming connections at the receiver’s side, as well as adding authentication data to outgoing requests at the exporter’s side.

This mechanism is implemented on top of the extensions framework provided by OTEL Collector.Authenticators are regular extensions that also satisfy one or more interfaces mentioned below :

**Server Authenticator Interface :**
(single interface for gRPC and HTTP)
<pre> 1. <a>go.opentelemetry.io/collector/config/configauth/ServerAuthenticator</a></pre>

**Client Authenticator Interfaces :**
<pre>
   1. <a>go.opentelemetry.io/collector/config/configauth/GRPCClientAuthenticator</a>
   2. <a>go.opentelemetry.io/collector/config/configauth/HTTPClientAuthenticator</a>
</pre>


**Typs of Authenticators :**

**1> Server Authenticator :**
 Used in receivers, intercept incoming HTTP/gRPC request (telemetry-signal), authentication data is expected to be contained in incoming telemetry-signals , this is used for authentication of telemetry signal. 

**2> Client Authenticator :**
 Used in exporters,  perform client-side authentication for outgoing telemetry signals, add authentication-data to telemetry-signal which is exported to downstream component e.g. Elastic-APM.

## 2.1 Server Authenticator : ##
It is an GOLang "extension" framework interface-implementation with a _**Authenticate()**_ funtion , receiving the payload-headers as parameter. 

If the authenticator is able to authenticate the incoming connection, then **return a nil error** , otherwise return concrete-error.

As an extension, the authenticator should make sure to initialize all the resources it needs during the **Start()** phase, and is expected to clean them up upon **Shutdown()**.

refer :
<pre><a>https://github.com/open-telemetry/opentelemetry-collector/blob/main/extension/auth/server.go</a>
<a>https://github.com/open-telemetry/opentelemetry-collector/blob/main/service/extensions/extensions.go</a></pre>

**The Authenticate call is part of the hot path for incoming requests and will block the pipeline if it does not complete execution.**
Thus it becomes critical to properly handle any blocking operations. Concretely, respect the deadline set by the context, in case one is provided. 

## 3.IMPLEMENTATION CONSIDERATIONS ##
  **1> Securing communication between OTEL collector and session service :**

      This is performed using one-way ssl ,session service sending server certificate and certificate being validated on otel-collector , ca-certificate would be used to validate SSL certificate provided by sesison-service , ca-certificate setting is optional , in case it is not set go will use system level settings for validating server certs (which might not work correctly) . This certificate is assumed to be mounted in a directory on the file system which is specified in collector run configuration (infa_auth.ca_cert_path)
  
  **2> Securing communication between agent and OTEL collector :**
  
      TO-DO
  
  **3> Securing communication between OTEL collector and Elastic APM server:**
  
      TO-DO 
  
  **4> NO session caching :**
  
     Any session info fetched including session expiry time , will not be cached , session service GET agent sesison API 
     would be invoked for each arriving telemetry-signal

## 4.IMPLEMENTATION DETAILS ##


## 5._BUILDING collector binary_ ##

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
##  _mock session service_ ##
rename package name to main and method-name to main and execute with go run , this api will be hosted :GET http://127.0.0.1:9898/session-service/api/v1/session/Agent ,
this API expects header IDS-AGENT-SESSION-ID : 123123123 to returns http 200 response , if header does not exist or value is different then API returns http 401 ,
I have not been able to figure out way to run this API as part of the test cases to this step is manual for now.
