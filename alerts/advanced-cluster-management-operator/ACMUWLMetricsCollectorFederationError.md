# ACMUWLMetricsCollectorFederationError

## Meaning

This alert fires when the ACM User Workload (UWL) Metrics Collector receives a high rate of non-2xx (e.g., 403 Forbidden, 500 Server Error) responses while trying to federate (scrape) metrics from the Hub cluster's User Workload Prometheus.

The query for this alert is `(sum by (status_code, type) (rate(acm_uwl_metrics_collector_federate_requests_total{status_code!~"2.*"}[10m]))) > 10`. This means the alert will fire if the rate of failed scrapes (with a valid HTTP error code) is greater than 10 per second for a continuous 10-minute period.

This component is different from the platform metrics collector. This one (`uwl-metrics-collector-deployment`) is specifically responsible for collecting metrics from your own applications in the `openshift-user-workload-monitoring` namespace.

## Impact

Metrics from user-defined workloads (e.g., custom application metrics) on the Hub cluster are not being collected by ACM Observability. This will cause user-created dashboards in Grafana to have no data.

**What user workload monitoring includes:**
* Custom Prometheus metrics exposed by your applications (e.g., `/metrics` endpoints)
* ServiceMonitors and PodMonitors created in user namespaces
* Custom application-specific metrics and alerts

**Who is affected:**
* Application developers relying on custom metrics dashboards
* Teams monitoring custom application SLIs/SLOs

This alert does not affect the collection of OpenShift platform metrics (cluster CPU, memory, etc.), which are handled by the separate `metrics-collector-deployment`.

## Diagnosis

The primary goal is to determine why the `uwl-metrics-collector-deployment` pod is receiving HTTP errors from the user workload Prometheus. The most likely cause is an RBAC (Role-Based Access Control) failure.

### 1. Verify User Workload Monitoring is enabled

Before investigating collector issues, confirm that User Workload Monitoring is enabled on your cluster:

```console
$ oc get configmap cluster-monitoring-config -n openshift-monitoring -o yaml | grep enableUserWorkload
```

If not present or set to `false`, user workload monitoring is not enabled. See the [OpenShift documentation](https://docs.openshift.com/container-platform/latest/monitoring/enabling-monitoring-for-user-defined-projects.html) for enabling it.

### 2. Check the logs of the uwl-metrics-collector-deployment pod

This is the most important step. The logs will show the exact error.

```console
$ oc logs -f deployment/uwl-metrics-collector-deployment -n open-cluster-management-observability
```

### 3. Analyze the log error

* **If you see `err="... i/o timeout"` or `err="... EOF"`**: 

    This is a network-level error, likely caused by a `NetworkPolicy`. This type of error will not trigger the alert, as it does not have an HTTP status code.

* **If you see `err="... 403 Forbidden"` or `err="... 401 Unauthorized"`**: 

    This is an authentication/authorization error. This is the cause of the alert. It means the collector is not authorized to scrape the user workload Prometheus.

### 4. If you see 403/401 errors, investigate the collector's RBAC permissions

* Find the `ServiceAccount` used by the `uwl-metrics-collector-deployment` pod:
    ```console
    $ oc get deployment uwl-metrics-collector-deployment -n open-cluster-management-observability -o jsonpath='{.spec.template.spec.serviceAccountName}'
    ```

    (This will likely be `uwl-metrics-collector-sa` or `endpoint-observability-operator-sa`)

* Find the `ClusterRoleBinding` that grants this `ServiceAccount` permission. (Replace `<SA_NAME>` with the name you found above):
    ```console
    $ oc get clusterrolebindings -o yaml | grep -B 10 -A 10 "<SA_NAME>"
    ```

* Look for the binding that grants this `ServiceAccount` the `cluster-monitoring-view` role or a specific role for `openshift-user-workload-monitoring`.

## Mitigation

The RBAC permissions for the `uwl-metrics-collector-deployment` are automatically managed by the `multicluster-observability-operator`. If the permissions are missing or incorrect, it means the operator's reconciliation loop has failed.

### 1. Check the logs of the multicluster-observability-operator

Check the logs of the `multicluster-observability-operator` (in the `open-cluster-management` namespace) for errors.

```console
$ oc logs -n open-cluster-management deployment/multicluster-observability-operator
```

Look for any errors related to "reconciling" or "failed to create ClusterRoleBinding".

### 2. Restart the multicluster-observability-operator pod

This is often the safest and fastest way to fix the issue. Restarting the operator will force it to re-run its reconciliation loop, and it will re-create or fix the missing/corrupted `ClusterRoleBinding` for the UWL collector.

```console
$ oc delete pod -n open-cluster-management -l name=multicluster-observability-operator
```

After the operator restarts, the `uwl-metrics-collector-deployment` pod's 403 errors should stop, and the alert will clear after the 10-minute window.
