# ACMUWLMetricsCollectorForwardRemoteWriteError

## Meaning

This alert fires when the ACM User Workload (UWL) Metrics Collector receives a high rate of non-2xx (e.g., 403 Forbidden, 503 Service Unavailable) responses while trying to "remote write" (forward) its scraped metrics to the central Observability API (`observability-observatorium-api`).

The alert's query is `(sum by (status_code, type) (rate(acm_uwl_metrics_collector_forward_write_requests_total{status_code!~"2.*"}[10m]))) > 10`. This means the alert will fire if the rate of failed remote-write requests (with a valid HTTP error code) is greater than 10 per second for a continuous 10-minute period.

This component (`uwl-metrics-collector-deployment`) is separate from the platform collector and is responsible only for metrics from your own applications (in `openshift-user-workload-monitoring`).

## Impact

User-defined workload metrics (e.g., custom application metrics) from the Hub cluster are not being successfully stored in the ACM Observability stack. This will cause custom Grafana dashboards to be empty or show "No Data".

This alert does not affect the collection of OpenShift platform metrics (cluster CPU, memory, etc.), which are handled by the separate `metrics-collector-deployment`.

## Diagnosis

The primary goal is to determine why the `uwl-metrics-collector-deployment` pod is failing to send data to the `observability-observatorium-api`.

### 1. Check the logs of the uwl-metrics-collector-deployment pod

This is the most important step. The logs will show the exact error.

```console
$ oc logs -f deployment/uwl-metrics-collector-deployment -n open-cluster-management-observability
```

### 2. Analyze the log error

* **If you see `err="... EOF"` or `err="... i/o timeout"`**: 

    This is a network-level error, not an HTTP error, and will not typically trigger this alert. This error means the collector cannot reach the `observability-observatorium-api` service at all. Check that the `observability-observatorium-api` pods are running.

* **If you see `err="... 401 Unauthorized"` or `err="... 403 Forbidden"`**: 

    This is an authentication/authorization error. This is the most likely cause of the alert. It means the `observability-observatorium-api` is rejecting the collector's connection.

### 3. If you see 401/403 errors, investigate the collector's mTLS certificate

The `uwl-metrics-collector-deployment` authenticates using an mTLS client certificate. This error means the certificate is likely missing, invalid, or has been rejected.

* Find the secret name for the pod's mTLS certificate:
    ```console
    $ oc get deployment uwl-metrics-collector-deployment -n open-cluster-management-observability -o yaml | grep secretName
    ```

* Look for a secret name related to client certs (e.g., `observability-controller-open-cluster-management.io-observability-signer-client-cert` or a `uwl-` specific version).

## Mitigation

This alert is almost always caused by a problem with the client certificate used for mTLS authentication. These certificates are managed by the `multicluster-observability-operator`.

### 1. If the diagnosis shows 401/403 errors (Authentication Failure)

The fastest and safest way to resolve this is to force the `multicluster-observability-operator` to reconcile its certificates. Restarting the operator will cause it to check all its managed components and reissue any missing or invalid certificates.

* Restart the `multicluster-observability-operator` pod:
    ```console
    $ oc delete pod -n open-cluster-management -l name=multicluster-observability-operator
    ```

* After the operator restarts, it will fix the secret. You must then restart the `uwl-metrics-collector-deployment` pod to pick up the new, valid secret.
    ```console
    $ oc scale deployment uwl-metrics-collector-deployment -n open-cluster-management-observability --replicas=0
    $ oc scale deployment uwl-metrics-collector-deployment -n open-cluster-management-observability --replicas=1
    ```

### 2. If the diagnosis shows EOF or i/o timeout (Network Failure)

This means the `observability-observatorium-api` service is down.

* Check the status of the `observability-observatorium-api` deployment and scale it up if needed:
    ```console
    $ oc get deployment observability-observatorium-api -n open-cluster-management-observability
    ```

* Scale the deployment if needed:
    ```console
    $ oc scale deployment observability-observatorium-api -n open-cluster-management-observability --replicas=2
    ```
