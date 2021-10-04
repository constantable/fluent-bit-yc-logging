![docker-publish](https://github.com/constantable/fluent-bit-yc-logging/actions/workflows/docker-publish.yml/badge.svg)
![go](https://github.com/constantable/fluent-bit-yc-logging/actions/workflows/go.yml/badge.svg)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/constantable/fluent-bit-yc-logging?style=flat)

# Yandex Cloud Logging Fluent Bit output plugin

[Fluent Bit](https://fluentbit.io) is a fast and lightweight log processor and forwarder or Linux, OSX and BSD family operating systems.

[Yandex Cloud Logging](https://cloud.yandex.com/en/docs/logging/) read and write logs of services and user applications by grouping messages into log groups.

## Installation

1. Create a Service account in Yandex Cloud
2. Create Authorized Key for the service account
3. Create Logging Group
4. Create your values.yaml as described in the [Configuration example](#-Configuration-example) section

To add the `fluent` helm repo, run:

```sh
helm repo add fluent https://fluent.github.io/helm-charts
```

To install a release named `fluent-bit`, run:

```sh
helm install fluent-bit fluent/fluent-bit --values example-values.yaml
```

## Configuration example

```yml
image:
  repository: ghcr.io/constantable/fluent-bit-yc-logging
  tag: "latest"
  pullPolicy: Always

config:
  service: |
    [SERVICE]
        Daemon Off
        Flush 1
        Log_Level info
        Parsers_File parsers.conf
        Parsers_File custom_parsers.conf
        HTTP_Server On
        HTTP_Listen 0.0.0.0
        HTTP_Port 2020
        Health_Check On

  inputs: |
    [INPUT]
        Name tail
        Path /var/log/containers/*_default_*.log
        multiline.parser docker, cri
        Tag kube.*
        Mem_Buf_Limit 5MB
        Skip_Long_Lines On
  filters: |
    [FILTER]
        Name kubernetes
        Match kube.*
        Merge_Log On
        Keep_Log Off
        K8S-Logging.Parser On
        K8S-Logging.Exclude On
  outputs: |
    [OUTPUT]
        Name fluent-bit-yc-logging
        Match kube.*
        ServiceAccountId <service account id>
        LogGroupId <log group id>
        KeyId <key id>
        PublicKey <base64 encoded public key>
        PrivateKey <base64 encoded private key>
        Retry_Limit 5

  customParsers: |
    [PARSER]
        Name docker_no_time
        Format json
        Time_Keep Off
        Time_Key time
        Time_Format %Y-%m-%dT%H:%M:%S.%L
```