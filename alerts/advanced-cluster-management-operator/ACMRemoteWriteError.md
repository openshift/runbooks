# ACMRemoteWriteError

## Meaning

This alert fires when the `observability-observatorium-api` service (on the Hub cluster) fails to send metrics to a configured external Prometheus remote-write endpoint.

The alert query `sum by (name) (rate(acm_remote_write_requests_total{code!~"2.*"}[10m])) / sum by (name) (rate(acm_remote_write_requests_total[10m])) > 0.2` means that the alert will fire if more than 20% of remote-write requests fail with a non-2xx HTTP status code (like 401 Unauthorized, 403 Forbidden, or 503 Service Unavailable) for a continuous 10-minute period.

## Impact

When this alert is firing, metrics from the ACM observability stack are not being successfully forwarded to the configured external monitoring system (e.g., a central Grafana, another Thanos instance, corporate monitoring platform, etc.).

**What is affected:**

* Metrics are not reaching the external system for centralized monitoring, alerting, or compliance reporting.
* Any dashboards, alerts, or SLI/SLO calculations in the external system that depend on ACM metrics will be incomplete or showing stale data.

**What is NOT affected:**

* Local ACM observability remains functional - you can still view metrics in the Hub cluster's Grafana.
* Metric collection from managed clusters continues normally.
* Data is still being stored in the local Thanos storage.

This does not affect the local ACM observability (you can still view metrics in the Hub cluster's Grafana). However, it signifies a broken integration and a loss of data on the external system.

## Diagnosis

The primary goal is to determine why the `observability-observatorium-api` pod is receiving HTTP errors from the external endpoint.

### 1. Check the logs of the observability-observatorium-api pods

This is the most important step. The logs will show the exact error code and the failing URL.

```console
oc logs -f deployment/observability-observatorium-api -n open-cluster-management-observability
```

Look for error messages like `level=error ... msg="failed to forward metrics" returncode="503 Service Temporarily Unavailable"`.

### 2. Inspect the MultiClusterObservability configuration

Check the `writeStorage` block to identify which secret is being used for the remote-write configuration.

```console
oc get mco observability -n open-cluster-management-observability -o yaml
```

Look for the `spec.storageConfig.writeStorage` section and the `name:` of the secret being used (e.g., `name: broken-remote-write`).

### 3. Inspect the remote-write Secret

Check the `ep.yaml` data within the secret to find the exact URL and credentials being used.

```console
oc get secret <YOUR_SECRET_NAME> -n open-cluster-management-observability -o jsonpath='{.data.ep\.yaml}' | base64 --decode
```

Example:

```yaml
url: http://httpbin.org/status/503
http_client_config:
  basic_auth:
    username: user
    password: wrong-password
```

### 4. Manually test the endpoint

Manually test the endpoint from a test pod (as seen in the KCS) to verify if the endpoint is reachable and if the credentials are correct.

**Common error patterns to look for:**

* `401 Unauthorized` or `403 Forbidden` - Authentication/authorization issue with credentials
* `503 Service Unavailable` or `502 Bad Gateway` - External service is down or overloaded
* `404 Not Found` - Incorrect URL or endpoint path
* `timeout` - Network connectivity issues or firewall blocking egress traffic

For more information about configuring remote-write endpoints, see [Exporting metrics to external endpoints](https://docs.redhat.com/en/documentation/red_hat_advanced_cluster_management_for_kubernetes/2.14/html-single/observability/index#exporting-metrics-to-external-endpoints).

## Mitigation

The mitigation depends on the errors found in the diagnosis.

### 1. If the url: in the secret is incorrect

Edit the secret and correct the `url:` to point to a valid remote-write endpoint.

```console
oc edit secret <YOUR_SECRET_NAME> -n open-cluster-management-observability
```

### 2. If the endpoint is returning 401/403 (Authentication Errors)

Edit the secret and correct the `http_client_config.basic_auth` (or other) credentials to match what the external service expects.

```console
oc edit secret <YOUR_SECRET_NAME> -n open-cluster-management-observability
```

### 3. If the endpoint is returning 5xx errors

This indicates the external remote-write service itself is unhealthy or unavailable. The issue must be resolved on the external system.

Check the external service's status and logs. In the meantime, you can monitor the ACM side to see when the service recovers:

```console
oc logs -f deployment/observability-observatorium-api -n open-cluster-management-observability
```

The alert will automatically clear once the external service becomes healthy and starts accepting requests again.

### 4. If the configuration is no longer needed

To stop the errors and clear the alert, you can disable the remote-write feature.

* Back up the MCO:

    ```console
    oc get mco observability -n open-cluster-management-observability -o yaml > mco-backup.yaml
    ```

* Edit the MCO:

    ```console
    oc edit mco observability -n open-cluster-management-observability
    ```

* Remove the entire `writeStorage:` block from the `spec:` section.

* Save the file. This will restart the `observability-observatorium-api` pods and stop them from trying to send metrics.
