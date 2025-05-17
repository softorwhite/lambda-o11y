# lambda-o11y


## What's This?

This repository contains example AWS Lambda code with OpenTelemetry tracing capability.
You can see the tracing results on AWS X-Ray, and this utilises [AWS Distro for OpenTelemetry Lambda's](https://aws-otel.github.io/docs/getting-started/lambda) Golang Lambda Layer which enables X-Ray integration with minimal settings.

## Prerequisites
- [Transaction Search](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Enable-Lambda-TransactionSearch.html) being enabled in CloudWatch
    - Set the Trace Indexing Rate high enough to get captured (100%)
  