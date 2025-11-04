# ACMMetricsCollectorForwardRemoteWriteError

## Meaning

This alert fires when the ACM Metrics Collector (running on the Hub cluster) fails to "remote write" its metrics to the central Thanos Hub (`observability-observatorium-api`) service.

The alert query `(sum by (status_code, type) (rate(acm_metrics_collector_forward_write_requests_total{status_code!~"2.*"}[10m]))) > 10` means the alert triggers only when the collector receives a high rate of non-2xx (e.g., 401, 403, 503) HTTP status codes for more than 10 minutes.

This alert will not fire for network-level failures like `i/o timeout` or `EOF` (connection dropped), as these do not produce an HTTP status code.

## Impact

When this alert is firing, metrics from the Hub cluster itself are successfully collected but cannot be sent to the central Thanos storage. This will cause the Hub cluster to appear "dark" or "offline" in Grafana, and its health and performance metrics will be missing. This can also negatively affect SLO/SLI calculations that rely on Hub cluster metrics.

## Diagnosis

The primary goal is to determine why the `metrics-collector-deployment` is receiving HTTP 4xx or 5xx errors.

### 1. Check the logs of the Metrics Collector pod

This is the most important step. The error message will tell you if the failure is a network error (like `EOF`) or an HTTP error (like 403).

```console
$ oc logs -f deployment/metrics-collector-deployment -n open-cluster-management-observability
```

### 2. Analyze the log error

* **If the log shows `err="... EOF"` or `err="... i/o timeout"`**:

    This is a network-level failure, which will not trigger this alert. This means the `observability-observatorium-api` service is likely down. Check that its pods are running:

    ```console
    $ oc get pods -n open-cluster-management-observability -l app.kubernetes.io/name=observatorium-api
    ```

    If no pods are Running, this is the root cause.

* **If the log shows `err="... response status 401 Unauthorized"` or `err="... 403 Forbidden"`**:

    This is an authentication error and will trigger the alert. This means the `observability-observatorium-api` server is running but is actively rejecting the collector's connection. This can be caused by:

    * The metrics-collector's client certificate (`observability-controller-open-cluster-management.io-observability-signer-client-cert` secret) is invalid or has expired.
    * The metrics-collector's `ServiceAccount` token is invalid or not authorized.
    * The `observability-observatorium-api` server's own mTLS configuration is broken and is failing to validate the (valid) client.

## Mitigation

The mitigation depends on the errors found in the diagnosis.

### 1. If observability-observatorium-api pods are not running (causing EOF errors)

The `observability-observatorium-api` deployment may have been scaled to 0. Scale the deployment back to its default replica count (usually 2).

```console
$ oc scale deployment observability-observatorium-api -n open-cluster-management-observability --replicas=2
```

### 2. If the collector is receiving 401/403 errors (the alert's trigger)

This indicates an authentication failure. The `metrics-collector` authenticates with both a client certificate and a `ServiceAccount` token.

* **Step 1**: Restart the `metrics-collector` pod. This will force it to reload its certificates and token.
    ```console
    $ oc scale deployment metrics-collector-deployment -n open-cluster-management-observability --replicas=0
    $ oc scale deployment metrics-collector-deployment -n open-cluster-management-observability --replicas=1
    ```

* **Step 2**: Restart the `observability-observatorium-api` pods. This will force the server to reload its authentication configuration.
    ```console
    $ oc scale deployment observability-observatorium-api -n open-cluster-management-observability --replicas=0
    $ oc scale deployment observability-observatorium-api -n open-cluster-management-observability --replicas=2
    ```

* **Step 3**: If errors persist, check the client certificate. The `observability-controller-open-cluster-management.io-observability-signer-client-cert` secret may be corrupted or expired. It is managed by the `multicluster-observability-operator`. Check the operator logs for errors:
    ```console
    $ oc logs -n open-cluster-management deployment/multicluster-observability-operator
    ```
