# ACMMetricsCollectorFederationError

## Meaning

This alert fires when the ACM Metrics Collector, which runs on the Hub cluster, receives a high rate of non-2xx (e.g., 401 Unauthorized, 403 Forbidden, 5xx Server Error) responses while trying to federate (scrape) metrics from the Hub cluster's internal Prometheus (`prometheus-k8s` in the `openshift-monitoring` namespace).

The alert expression `(sum by (status_code, type) (rate(acm_metrics_collector_federate_requests_total{status_code!~"2.*"}[10m]))) > 10` indicates that the rate of failed requests has exceeded 10 per second for 10 minutes (600s).

## Impact

When this alert is firing, metrics from the Hub cluster itself are not being collected by ACM Observability. This will cause the Hub cluster to appear "dark" or "offline" in Grafana dashboards, as its own performance and health data will be missing. This can also negatively affect SLO/SLI calculations that rely on Hub cluster metrics.

## Diagnosis

### 1. Check the logs of the Metrics Collector pod

The `metrics-collector-deployment` pod in the `open-cluster-management-observability` namespace is responsible for this scrape. Inspect its logs for specific errors:

```console
$ oc logs -f deployment/metrics-collector-deployment -n open-cluster-management-observability
```

Look for repeated HTTP 401/403 (Authorization) or 5xx (Server) errors.

* **`i/o timeout` errors** indicate a network problem, not an authorization problem.

* **401 Unauthorized or 403 Forbidden errors** indicate an RBAC problem.

### 2. Verify the ClusterRoleBinding for the collector

The `metrics-collector-deployment` pod uses a `ServiceAccount` that requires permission to get metrics from the `openshift-monitoring` namespace. Verify this binding is correct.

* Find the `ServiceAccount`:
    ```console
    $ oc get deployment metrics-collector-deployment -n open-cluster-management-observability -o jsonpath='{.spec.template.spec.serviceAccountName}'
    endpoint-observability-operator-sa
    ```

* Find its `ClusterRoleBinding`:
    ```console
    $ oc get clusterrolebinding metrics-collector-view
    NAME                     ROLE                            AGE
    metrics-collector-view   ClusterRole/cluster-monitoring-view   ...
    ```

* Inspect the `ClusterRoleBinding`:
    ```console
    $ oc get clusterrolebinding metrics-collector-view -o yaml
    ```

* Ensure the `roleRef.name` is `cluster-monitoring-view` and the `subjects` section correctly lists the `endpoint-observability-operator-sa` `ServiceAccount` from the `open-cluster-management-observability` namespace.

### 3. Check for NetworkPolicy blocking traffic

If you see `i/o timeout` or connection refused errors, check for NetworkPolicies that might be blocking traffic:

```console
$ oc get networkpolicy -n open-cluster-management-observability
$ oc get networkpolicy -n openshift-monitoring
```

Look for policies that might restrict egress from the metrics-collector pod or ingress to the prometheus-k8s service.


## Mitigation

Resolution involves restoring the `metrics-collector-deployment`'s ability to successfully authenticate and scrape metrics from the `openshift-monitoring` namespace.

### 1. If the ClusterRoleBinding is missing or incorrect

The `metrics-collector-view` `ClusterRoleBinding` is managed by the `multicluster-observability-operator`. If it has been deleted or modified, it will not be automatically recreated if the operator is paused or in a failed state.

* Ensure the `multicluster-observability-operator` and `multiclusterhub-operator` in the `open-cluster-management` namespace are running.

* If the binding is missing, you may need to force a reconciliation by restarting the `multicluster-observability-operator` pod.

* If the binding cannot be restored automatically, it can be reapplied from a backup.

### 2. If a NetworkPolicy is causing the issue

If a custom `NetworkPolicy` was recently applied, it may be blocking the `metrics-collector-deployment` pod from reaching the `prometheus-k8s` service on port 9091 in the `openshift-monitoring` namespace.

* Remove or modify the `NetworkPolicy` to allow this traffic.

### 3. If openshift-monitoring components are failing

If the collector logs show 5xx errors, the problem may be with Prometheus itself.

* Check the health of the `prometheus-k8s` pods in the `openshift-monitoring` namespace:
    ```console
    $ oc get pods -n openshift-monitoring -l app.kubernetes.io/name=prometheus
    ```

* Check the logs of the `kube-rbac-proxy` container within the `prometheus-k8s` pods, as it may be failing its own authentication checks against the Kube API server:
    ```console
    $ oc logs -n openshift-monitoring prometheus-k8s-0 -c kube-rbac-proxy
    ```

You should see successful scrape operations with no 403 errors. The alert will clear after 10+ minutes of successful federation.

**Expected resolution time**: Alert should clear within 10-15 minutes after RBAC permissions are restored and successful scraping resumes.
