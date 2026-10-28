# MetricsServerKubeletScrapeFailures

## Meaning

The `MetricsServerKubeletScrapeFailures` alert fires when a `metrics-server`
pod fails at least 10% of kubelet scrape requests over the last 5 minutes
for at least 15 minutes. The alert has `warning` severity.

## Impact

`metrics-server` collects CPU and memory usage from kubelets and serves them
through the Resource Metrics API (`metrics.k8s.io`). This API powers
`oc adm top`, the Horizontal Pod Autoscaler (HPA), and the Vertical Pod
Autoscaler (VPA).

When this alert fires:

- `oc adm top nodes` and `oc adm top pods` may report `<unknown>` for affected
  nodes or pods.
- HPA and VPA may be unable to compute resource utilization and fail to scale
  workloads correctly.
- Cluster administrators lose visibility into live resource usage through the
  Kubernetes Resource Metrics API.

This alert does not affect the Prometheus monitoring pipeline. Prometheus
collects resource metrics independently from cAdvisor/kubelet scrapes.

## Diagnosis

The alert includes `namespace` and `pod` labels identifying the affected
metrics-server instance. It does not include a `node` label; identify
affected nodes from the metrics-server logs or PromQL queries below.

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

Review the graph to understand whether the issue is new or if there were prior
signs of elevated failure rates.

To identify nodes for which metrics-server experiences high scrape latency:

```console
histogram_quantile(0.99, rate(metrics_server_kubelet_request_duration_seconds_bucket[15m]))
```

To identify nodes that exist but have not been scraped recently by
metrics-server:

```console
(time() - metrics_server_kubelet_last_request_time_seconds) and on (node) max by(node) (kube_node_info) > 60
```

### Review metrics-server logs

Check the metrics-server deployment logs for kubelet scrape errors:

```console
$ oc -n openshift-monitoring logs pod/<pod-from-alert> -c metrics-server --tail=200
```

Increase `--tail` for larger clusters where many nodes may be failing scrapes.

Common log patterns include (exact wording may vary by metrics-server version):

- `Failed to scrape node, node is not ready` — the node exists but is not in a
  Ready state. Nodes in `NotReady` state are expected to fail scrapes and
  contribute to the failure ratio. This is normal during node lifecycle events
  such as upgrades, reboots, or scaling. Note that this message may not be
  visible at the configured log verbosity.
- `Failed to scrape node, timeout to access kubelet` — metrics-server could not
  reach the kubelet within the configured scrape timeout.
- Other `Failed to scrape node` errors — often TLS, network, or kubelet
  availability problems.

Each log line includes the affected `node` name.

### Check Resource Metrics API health

Verify whether resource metrics are missing from the API:

```console
$ oc adm top nodes
```

Nodes reporting `<unknown>` are likely failing kubelet scrapes.

You can also query the Resource Metrics API directly:

```console
$ oc get --raw /apis/metrics.k8s.io/v1beta1/nodes
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

If specific nodes consistently fail scrapes:

- Review kubelet logs on the affected node. See the [KubeletDown](KubeletDown.md)
  runbook for guidance on accessing node and kubelet logs.
- Check for resource pressure, disk issues, or network problems on the node.

Test whether the kubelet exposes resource metrics:

```console
$ NODE_NAME=<affected node>
$ oc get --raw /api/v1/nodes/$NODE_NAME/proxy/metrics/resource
```

This checks access to the kubelet resource metrics endpoint through the API server.
Investigate any returned error in that context. A successful response does not verify
metrics-server's direct connectivity or client-certificate configuration.

### Check metrics-server configuration and certificates

If scrape failures affect many or all nodes:

- Verify the `metrics-server-client-certs` secret and `kubelet-serving-ca-bundle`
  config map exist in `openshift-monitoring`:

  ```console
  $ oc -n openshift-monitoring get secret metrics-server-client-certs
  $ oc -n openshift-monitoring get configmap kubelet-serving-ca-bundle
  ```

- Check whether the metrics-server deployment was recently modified or is
  failing to roll out:

  ```console
  $ oc -n openshift-monitoring rollout status deployment/metrics-server
  ```

- Review cluster-monitoring-operator logs if the deployment is not healthy:

  ```console
  $ oc -n openshift-monitoring logs deployment/cluster-monitoring-operator
  ```

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

### Network connectivity

metrics-server must reach each kubelet on the kubelet's serving port (10250 by
default). Firewalls, routing issues, or network policies affecting
metrics-server pods can cause scrape timeouts.

Review cluster network health and any recent changes to firewalls or network
policies in the `openshift-monitoring` namespace.

### Persistent or cluster-wide failures

If the alert does not resolve after addressing node-level issues, gather
metrics-server logs, node status, and the PromQL query results above and open a
support case.
