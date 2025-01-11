# PrometheusKubernetesListWatchFailures

## Meaning

The `PrometheusKubernetesListWatchFailures` alert is triggered when there is a constant
increase in failures with `LIST/WATCH` requests to the Kubernetes API during Prometheus target
discovery.

## Impact

This may prevent Prometheus from adding, updating, or deleting targets effectively.

## Diagnosis

Determine whether the alert has triggered for the instance of Prometheus used
for default cluster monitoring or for the instance that monitors user-defined
projects by viewing the alert message's `namespace` label: the namespace for
default cluster monitoring is `openshift-monitoring` and the namespace for
user workload monitoring is `openshift-user-workload-monitoring`.

## Mitigation

To gain further insight, review the logs of the affected Prometheus instance:

   ```shell
   $ NAMESPACE='<value of namespace label from alert>'

   $ oc -n $NAMESPACE logs -c prometheus -l 'app.kubernetes.io/name=prometheus'
   ```

The issue may arise from one or both of the following scenarios:

### Insufficient RBAC Permissions for Prometheus

In this scenario, Prometheus has been tasked to discover targets in a specified namespace, likely
through `ServiceMonitor` or `PodMonitor` resources, yet Prometheus lacks the necessary RBAC
permissions to query `Service`, `Endpoints`, `Pod`, and other related resources where the targets 
are defined, the following log messages may be observed:

```
ts=2024-11-13T07:09:51.190Z caller=klog.go:108 level=warn component=k8s_client_runtime
func=Warningf msg="github.com/prometheus/prometheus/discovery/kubernetes/kubernetes.go:554:
failed to list *v1.Endpoints: endpoints is forbidden:
User \"system:serviceaccount:openshift-monitoring:prometheus-k8s\" cannot list resource
\"endpoints\" in API group \"\" in the namespace \"foo\""
```

To rectify this issue, ensure Prometheus is granted the necessary RBAC permissions as detailed in
the following guide: [Configuring Prometheus to scrape metrics].

---

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.

[Configuring Prometheus to scrape metrics]: https://rhobs-handbook.netlify.app/products/openshiftmonitoring/collecting_metrics.md/#configuring-prometheus-to-scrape-metrics

