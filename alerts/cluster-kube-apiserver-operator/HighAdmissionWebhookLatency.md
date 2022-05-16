# HighAdmissionWebhookLatency

## Meaning

An admission webhook is adding a significant amount of latency to API requests.

## Impact

Webhook is slowing down API requests. Performance degradation.

## Diagnosis

The alert should identify the name of the webhook.

> **NOTE:** There might be more than one webhook with the same name.

1. Find the webhook configuration.

```console
$ oc something...
```

2. Identify service.
3. Identify service workload and examine the workload logs.
4. Examine apiserver logs. 

## Mitigation

Fix problem with webhook.
