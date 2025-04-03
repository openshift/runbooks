# PrometheusKubernetesListWatchFailures

## Meaning

The `PrometheusKubernetesListWatchFailures` alert is triggered when `LIST/WATCH`
requests to the Kubernetes API start failing during Prometheus target discovery.

## Impact

This might prevent Prometheus from adding, updating, or deleting targets effectively.

## Diagnosis

Determine if the alert has triggered for the instance of Prometheus used
for default cluster monitoring or for the instance that monitors user-defined
projects by viewing the alert message's `namespace` label. The namespace for
default cluster monitoring is `openshift-monitoring` and the namespace for
user workload monitoring is `openshift-user-workload-monitoring`.

To gather more information, review the logs of the affected Prometheus instance:

   ```shell
   # NAMESPACE='<value of namespace label from alert>'
   $ oc -n $NAMESPACE logs -c prometheus -l 'app.kubernetes.io/name=prometheus' --tail=-1
   ```

## Mitigation

The issue might arise for multiple reasons:

### Insufficient RBAC Permissions for Prometheus

In this scenario, Prometheus has been tasked to discover targets in a specified
namespace, typically through `ServiceMonitor` or `PodMonitor` resources. However,
Prometheus does not have the necessary RBAC permissions to query `Service`, `Endpoints`,
`Pod`, and other related resources where the targets are defined. Logs similar to
the following one can be observed:

```text
... failed to list *v1.Endpoints: endpoints is forbidden:
User \"system:serviceaccount:openshift-monitoring:prometheus-k8s\" cannot list resource
\"endpoints\" in API group \"\" in the namespace \"foo\""
```

To resolve this issue, follow one of these steps:

- For default cluster monitoring (`openshift-monitoring`), follow
[Configuring Prometheus to scrape metrics] to ensure Prometheus is granted
the necessary RBAC permissions.

- For user workload monitoring (`openshift-user-workload-monitoring`), typically,
Prometheus deployed for user workload monitoring has cluster-wide permissions. Ensure
that cluster-monitoring-operator (CMO) has no issues while deploying the underlying
resources to obtain these permissions. Running the
`oc get co monitoring -o yaml` command should provide hints.

### Connectivity Issue

Prometheus might not be able to reach the API server. In such cases, logs similar
to the following ones can be seen (examples are not exhaustive):

```text
... failed to list *v1.Endpoints: Get
\"https://A.B.C.D:443/api/v1/namespaces/openshift-foo/endpoints?limit=500&resourceVersion=0\":
dial tcp A.B.C.D:443: i/o timeout"
```

```text
... failed to list *v1.Endpoints: Get
\"https://foo-apiserver:6443/api/v1/namespaces/openshift-console/endpoints?limit=500&
resourceVersion=0\": dial tcp: lookup foo-apiserver on A.B.C.D:53: no such host"
```

---

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.

[Configuring Prometheus to scrape metrics]: https://rhobs-handbook.netlify.app/products/openshiftmonitoring/collecting_metrics.md/#configuring-prometheus-to-scrape-metrics
