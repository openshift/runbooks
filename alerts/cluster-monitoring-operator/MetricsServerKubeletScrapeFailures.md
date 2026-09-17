# MetricsServerKubeletScrapeFailures

## Meaning

The `MetricsServerKubeletScrapeFailures` alert fires when a metrics-server pod
in the `openshift-monitoring` namespace fails more than 10% of its kubelet
scrape requests over a 15 minute period.

metrics-server collects CPU and memory usage from kubelets and serves them
through the Resource Metrics API (`metrics.k8s.io`). This API powers
`oc top`, the Horizontal Pod Autoscaler (HPA), and the Vertical Pod
Autoscaler (VPA).

A small background rate of kubelet scrape failures (typically below 2%) is
normal on busy or autoscaling clusters. The alert threshold is set at 10% to
avoid false positives.

## Impact

When this alert fires:

- `oc top nodes` and `oc top pods` may report `<unknown>` for affected nodes or
  pods.
- HPA and VPA may be unable to compute resource utilization and fail to scale
  workloads correctly.
- Cluster operators lose visibility into live resource usage through the
  Kubernetes metrics API.

This alert does not affect the Prometheus monitoring pipeline. Prometheus
collects resource metrics independently from cAdvisor/kubelet scrapes.

## Diagnosis

The alert includes `namespace` and `pod` labels identifying the affected
metrics-server instance.

### Review metrics-server logs

Check the metrics-server deployment logs for kubelet scrape errors:

```console
$ oc -n openshift-monitoring logs deployment/metrics-server --tail=200 | grep "Failed to scrape node"
```

Common log messages include:

- `Failed to scrape node, node is not ready` — the node exists but is not in a
  Ready state.
- `Failed to scrape node, timeout to access kubelet` — metrics-server could not
  reach the kubelet within the configured scrape timeout.
- Other `Failed to scrape node` errors — often TLS, network, or kubelet
  availability problems.

Each log line includes the affected `node` name.

### Check Resource Metrics API health

Verify whether resource metrics are missing from the API:

```console
$ oc top nodes
```

Nodes reporting `<unknown>` are likely failing kubelet scrapes.

You can also query the Resource Metrics API directly:

```console
$ oc get --raw /apis/metrics.k8s.io/v1beta1/nodes
```

### Inspect scrape failure rate

Open **Observe** -> **Metrics** in the OpenShift web console and run the
following PromQL expression to see the current failure ratio:

```console
sum by (pod, namespace) (
  rate(metrics_server_kubelet_request_total{namespace="openshift-monitoring",success="false"}[5m])
)
/
sum by (pod, namespace) (
  rate(metrics_server_kubelet_request_total{namespace="openshift-monitoring"}[5m])
)
```

### Check metrics-server pod health

```console
$ oc -n openshift-monitoring get deployment metrics-server
$ oc -n openshift-monitoring get pods -l app.kubernetes.io/name=metrics-server
$ oc -n openshift-monitoring describe deployment metrics-server
```

Confirm the pods are Ready and not restarting frequently.

### Check affected nodes

For each node reported in the metrics-server logs:

```console
$ oc get node $NODE_NAME
$ oc describe node $NODE_NAME
```

Review node conditions, recent events, and kubelet health.

Test whether the kubelet exposes resource metrics:

```console
$ NODE_NAME=<affected node>
$ oc get --raw /api/v1/nodes/$NODE_NAME/proxy/metrics/resource
```

If this request fails, times out, or returns zero usage values, the problem is
likely on the node or kubelet side rather than in metrics-server itself.

## Mitigation

The resolution depends on the root cause identified during diagnosis.

### Transient node conditions

Scrape failures often increase during normal cluster operations:

- Node upgrades or reboots
- Nodes being drained or cordoned
- Spot instance interruptions or autoscaling node churn
- Nodes in `NotReady` state during startup or recovery

If failures correlate with known node lifecycle events and resolve as nodes
return to Ready, no further action may be required.

### Kubelet or node issues

If specific nodes consistently fail scrapes:

- Review kubelet logs on the affected node. See the [KubeletDown](KubeletDown.md)
  runbook for guidance on accessing node and kubelet logs.
- Check for resource pressure, disk issues, or network problems on the node.
- Verify the kubelet `/metrics/resource` endpoint responds as described in the
  diagnosis section.

### Network connectivity

metrics-server must reach each kubelet on the kubelet's serving port (10250 by
default). Network policies, firewalls, or routing issues between metrics-server
pods and worker nodes can cause scrape timeouts.

Review cluster network health and any recent changes to network policies in the
`openshift-monitoring` namespace or on worker nodes.

### metrics-server configuration or certificates

If scrape failures affect many or all nodes:

- Verify the `metrics-server-client-certs` and `kubelet-serving-ca-bundle`
  secrets exist in `openshift-monitoring`:

  ```console
  $ oc -n openshift-monitoring get secret metrics-server-client-certs
  $ oc -n openshift-monitoring get configmap kubelet-serving-ca-bundle
  ```

- Check whether the metrics-server deployment was recently modified or is
  failing to roll out:

  ```console
  $ oc -n openshift-monitoring rollout status deployment/metrics-server
  ```

- Review cluster-monitoring-operator and metrics-server operator logs if the
  deployment is not healthy:

  ```console
  $ oc -n openshift-monitoring logs deployment/cluster-monitoring-operator
  ```

### Persistent or cluster-wide failures

If the alert does not resolve after addressing node-level issues, gather
metrics-server logs, node status, and the PromQL query results above and open a
support case.
