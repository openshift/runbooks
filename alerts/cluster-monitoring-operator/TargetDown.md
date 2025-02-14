# TargetDown

This runbook provides guidance for diagnosing and resolving the `TargetDown` alert
in OpenShift Container Platform.

## Meaning

The `TargetDown` alert fires when Prometheus has been unable to scrape one
or more targets over a specific period of time. It is triggered when specific
scrape targets within a service remain unreachable (`up` metric = 0) for a
predetermined duration.

## Impact

- **Visibility**: If a target is down, the metrics from the affected targets will
  not be captured by Prometheus. If metrics are not captured, you will have only
  limited insights about the health and performance of the associated application.
- **Alerts**: If a target is down, the accuracy of certain alerts can be compromised.
  For example, critical alerts might not be triggered, potentially causing service
  disruptions to go undetected.
- **Resource Optimization**: Auto-scalers might not function correctly if essential
  metrics are missing, which can result in wasted resources or a degraded user
  experience for your applications.

## Diagnosis and mitigation

### Identifying targets that are down

- Navigate to **Observe** -> **Targets** in the OpenShift web console. Choose **Down**
  in the **Filter** combo button to quickly list down targets. Click on individual
  targets for details and error messages from the last scrape attempt.
- Alternatively, query `up == 0` in **Observe** -> **Metrics** in the OpenShift web
  console. The metric labels will help pinpoint the affected Prometheus instance
  and the down target.

The alert and the metric `up` have the `namespace`, `service`， `job`, `pod`
and `prometheus` labels.
With these labels, we can identify which Prometheus instance fails to scrape which
target:

- The label `prometheus` should be either `openshift-monitoring/k8s` or `openshift-user-workload-monitoring/user-workload`,
  indicating the Prometheus pods scraping the target are either `prometheus-k8s-0`
  / `prometheus-k8s-1` in the namespace `openshift-monitoring`, or `prometheus-user-workload-0`
  / `prometheus-user-workload-1` in the namespace `openshift-user-workload-monitoring`.
- The label `pod` indicates the pod exposing the metrics endpoint.
- The labels `namespace` and `service` that help us locate the `Service` exposing
  the metric endpoint.
- The label `namespace` and `job` can locate the `ServiceMonitor` or `PodMonitor`
  that configures Prometheus to scrape the target. The `job` label is the name of
  the monitor.

Now we have both end of the metric scraping flow as well as the monitor resources
linking them together.
We are ready to diagnose the root cause by inspecting each component on the scraping
workflow.

### Potential Issues and Resolutions
#### Network Issues

Check the network connectivity between the Prometheus pod and the target pod.
Ensure that the target pod is reachable from the Prometheus pod and that there are
no firewall rules or network interruptions blocking communication.

There are some useful metrics to help investigating network issues:
- `net_conntrack_dialer_conn_attempted_total`
- `net_conntrack_dialer_conn_closed_total`
- `net_conntrack_dialer_conn_established_total`
- `net_conntrack_dialer_conn_failed_total`

OpenShift guide on [Troubleshooting network issues](https://docs.openshift.com/container-platform/4.17/support/troubleshooting/troubleshooting-network-issues.html)
provides more details.

#### Target Resource Exhaustion

Check if the pod's health and ready probes reports a good state.
Then check whether the metric endpoint is responsive.
We can either forward the metric port to local and then query it using use `curl`,
 `wget` or similar tools. Or on the container exposing the metric, send
query to the port serving the metric endpoint.

If the query returns an error, it is probably an application problem.
If the query takes too long or times out, it is probably due to resource exhaustion.

Some applications may enforce rate limiting or throttling. Check if the scraping
traffic is hitting such limits, causing the target to become temporarily unavailable.
Checking the pod's logs and events may help us diagnose such problems.

Then we should check if the target pod resource utilization is too high, causing
it to become unresponsive.
We can refer to the tab **Observe** -> **Dashboards** in the Openshift web Console
to have an overview of resource utilization. The Dashboard `Kubernetes/Compute Resources/Pod`
can show the CPU, memory, network and storage usage of the pod.

To view more details of CPU and memory usage, we can check these metrics:
- CPU Usage
  * `container_cpu_usage_seconds_total`: Cumulative CPU usage in seconds.
  * `pod:container_cpu_usage:sum` Represents the average CPU load of a pod over
  the last 5 minutes.
  * `container_cpu_cfs_throttled_seconds_total` Duration when the CPU was throttled
  due to limits.
- Memory Usage
  * `container_memory_usage_bytes` Current memory usage in bytes.
  * `container_memory_rss` Resident set size in bytes. This metric gives an
  understanding of the memory actively used by a pod.
  * `container_memory_swap` Swap memory usage in bytes.
  * `container_memory_working_set_bytes` Total memory in use. This includes all
  memory regardless of when it was accessed.
  * `container_memory_failcnt` Cumulative count of memory allocation failures.
- File System Usage
  * `container_fs_reads_bytes_total` and `container_fs_writes_bytes_total`
  Cumulative filesystem read and write operations in bytes.
  * `container_fs_reads_total` and `container_fs_writes_total` Cumulative
  filesystem read and write operations in bytes.
  * `kubelet_volume_stats_available_bytes` The free space in a volume in bytes.

If some basic metrics are not available, we can also use this command to get
 CPU and memory usage of a pod:
```bash
oc top pod $POD_NAME
```
As well as this command for volume free space:
```bash
oc exec -n $NAMESPACE $POD_NAME -- df
```

#### Target Application or Service Failure

Investigate if the application or service running on the target is experiencing
issues or has even crashed. Review logs and last metrics from the target pod to
identify any errors or crashes.

Here is [a guide to investigate pod issues](https://docs.openshift.com/container-platform/4.17/support/troubleshooting/investigating-pod-issues.html)
in OpenShift.

#### Incorrect Target Configuration

Verify that the scrape target configuration in Prometheus is correct. This
configuration is generated from the `ServiceMonitor` and `PodMonitors`.
We can get the `ServiceMonitor` and `PodMonitors` related to the alert with this
command using the `namespace` and `job` label. The name of `ServiceMonitor` is the
`job` name unless the property `jobLabel` is set in `ServiceMonitor`.
```bash
oc get servicemonitor $SERVICE_MONITOR_NAME -n $NAMESPACE -o yaml
```
Check the target's port, selector, scheme, TLS settings, etc for invalid values.
Please refer to [the Prometheus Operator API document](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md)
for detailed specification of `ServiceMonitor` and `PodMonitors`.

About the TLS settings, even if the values are correct, it may still fail due to
certificate expiration. Please refer to the next section for details.

#### Expired SSL/TLS Certificates

If the target uses SSL/TLS for communication, check if the SSL/TLS certificate
has expired or certificate files accessible by Prometheus.

Prometheus use a certificate to scrape the metrics endpoint is indicated in the
property `.spec.endpoints.tlsConfig.certFile` of a `ServiceMonitor` or `PodMonitor`.
The path of the certificate file points to a mounted volume on the Prometheus Pod.
Therefore, we can deduce which secret holds the certificate.

Here is an example.
We have a `ServiceMonitor` using `/etc/prometheus/secrets/metrics-client-certs/tls.crt`
as its certificate file.
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: exporter
    app.kubernetes.io/name: node-exporter
    app.kubernetes.io/part-of: openshift-monitoring
    app.kubernetes.io/version: 1.5.0
    monitoring.openshift.io/collection-profile: full
  name: node-exporter
  namespace: openshift-monitoring
spec:
  endpoints:
  - bearerTokenSecret:
      key: ""
    interval: 15s
    port: https
    relabelings:
    - action: replace
      regex: (.*)
      replacement: $1
      sourceLabels:
      - __meta_kubernetes_pod_node_name
      targetLabel: instance
    scheme: https
    tlsConfig:
      ca: {}
      caFile: /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt
      cert: {}
      certFile: /etc/prometheus/secrets/metrics-client-certs/tls.crt
      keyFile: /etc/prometheus/secrets/metrics-client-certs/tls.key
      serverName: node-exporter.openshift-monitoring.svc
  jobLabel: app.kubernetes.io/name
  namespaceSelector: {}
  selector:
    matchLabels:
      app.kubernetes.io/component: exporter
      app.kubernetes.io/name: node-exporter
      app.kubernetes.io/part-of: openshift-monitoring
```
The `Prometheus` that scrapes this target has a volume mounted at `/etc/prometheus/secrets/*`
where the secrets have the same name as the subdirectory. In this example, the secret
is `metrics-client-certs` in the namespace `openshift-monitoring`. Now we extract
the certificate from the secret using this command:
```bash
oc extract secret/$SECRET_NAME -n $NAMESPACE --keys=tls.crt --to=- > certificate.crt
```
Then we inspect its expiration date.
```bash
openssl x509 -noout -enddate -in certificate.crt
```
The output should contain the `notAfter` field as its expiration date.
```bash
notAfter=Aug  6 13:11:20 2025 GMT
```

Openshift usually **renews certificates automatically** before the expiration date.
This date should be sometime in the future. If the certificate does expire without
automatic renewal, please contact the OpenShift support team.

If the issue requires immediate resolution, please refer to [this guide on how to
force a certificate renewal](https://access.redhat.com/solutions/5899121).

To diagnose potential issues with automatic certificate renewal, perform the
following checks:
1. Ensure the `prometheus-k8s` stateful set logs do not display errors related to
   certificates.
   ```bash
   oc logs statefulset/prometheus-k8s -n openshift-monitoring
   ```
2. Verify that the scraped target pod(s) logs are free from certificate-related errors.
3. Check that the `cluster-monitoring-operator` pod is running and its logs contain
   no errors regarding certificate rotation.
   ```bash
   oc logs deployment/cluster-monitoring-operator -n openshift-monitoring
   ```
4. Confirm the `prometheus-operator` pod is operational and its logs have no errors
   concerning certificate rotation.
   ```bash
   oc logs deployment/prometheus-operator -n openshift-monitoring
   ```

#### Prometheus Scraping Interval

Ensure that the scraping interval for the target is appropriately configured.
Too frequent scraping can degrade the performance of target pod, even causing
timeout when scraping.

The last scrape duration is visible on the **Observe** -> **Targets** tab in the
OpenShift web Console.

Otherwise, we can check the metric `scrape_duration_seconds`.

If scrape duration is close to the scrape interval, we may consider to increase
the interval.

The scrape interval is configured in the `ServiceMonitor` or `PodMonitor`.


#### Target metrics path

Verify that the metrics path in the target is correctly defined in the `ServiceMonitor`
and `PodMonitor` resources. If the scraped target pod changes its metrics path in
recent code update, the metric endpoint should update accordingly.

#### Target restart or update

If the target was recently restarted or updated without replications, it might
temporarily become unavailable. Check if this is the case and allow some time for
it to stabilize.

#### Node failures

Check if the node where the target is running has experienced failures or evictions.
See [Verifying Node Health](https://docs.openshift.com/container-platform/4.17/support/troubleshooting/verifying-node-health.html)
for more information.
