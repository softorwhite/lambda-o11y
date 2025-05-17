# lambda-o11y


## What's This?

This repository contains example AWS Lambda code with OpenTelemetry tracing capability.
You can see the tracing results on AWS X-Ray, and this utilises [AWS Distro for OpenTelemetry Lambda's](https://aws-otel.github.io/docs/getting-started/lambda) Golang Lambda Layer which enables X-Ray integration with minimal settings.

## Prerequisites
- [Transaction Search](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Enable-Lambda-TransactionSearch.html) being enabled in CloudWatch
    - Set the Trace Indexing Rate high enough to get captured (100%)

## Expected Results
- Tracing data will show up on X-Ray console. <img width="1103" alt="スクリーンショット 2025-05-17 21 43 45" src="https://github.com/user-attachments/assets/6ae661d6-8d4e-4c27-91b6-324a0de7a868" />
- Events added by span.AddEvent() function **WILL NOT** appear in the AWS console since the event functionality is not compatible with X-Ray.
