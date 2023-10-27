#source image
#https://infacloud.jfrog.io/artifactory/ct-docker/golang/1.20/

FROM infacloud-ct-docker.jfrog.io/golang:1.20

RUN mkdir /opt/collector
WORKDIR /opt/collector

COPY target/otelcol-custom /opt/collector/
#COPY /testdata/otel_config.yaml /opt/collector/
COPY infa_auth/testdata/scheduler-service-keystore.jks /opt/collector/
COPY infa_auth/testdata/truststore.jks /opt/collector/

ENTRYPOINT ["./otelcol-custom" ,"--config", "/opt/resources/otel_config.yaml"]
EXPOSE 55681
EXPOSE 55680
