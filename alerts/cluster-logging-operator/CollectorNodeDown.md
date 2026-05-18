# CollectorNodeDown

## Meaning

The `CollectorNodeDown` alert indicates that Prometheus is unable to scrape
metrics from a log collector pod for more than 10 minutes. This means the
collector's metrics endpoint is unreachable, which could indicate that the
collector pod is down, has crashed, is experiencing network connectivity issues,
or its metrics port is not responding.

## Impact

When this alert fires, the following impacts may occur:

- **Loss of observability:** Metrics for log collection, forwarding rates, and
error rates are unavailable, making it impossible to monitor collector health
and performance.
- **Potential log data loss:** If the collector is actually down (not just
unreachable by Prometheus), logs are not being collected or forwarded to their
destinations.
- **Inability to detect other issues:** Without metrics, other collector-related
alerts (such as `ClusterLogForwarderOutputErrorRate`, `DiskBufferUsage`, or
`CollectorHigh403ForbiddenResponseRate`) cannot fire, masking additional problems.

## Diagnosis

1. **Check if the collector pod is running:**

   ```bash
   oc get pods -n openshift-logging -l app.kubernetes.io/component=collector
   ```

   Look for pods in a non-Running state (e.g., `CrashLoopBackOff`, `Error`,
   `Pending`).

2. **Examine pod status and recent events:**

   ```bash
   oc describe pod <POD_NAME> -n openshift-logging
   ```

   Check the "Events" section for errors such as:
   - Image pull failures
   - Resource constraints (CPU, memory)
   - Volume mount issues
   - Liveness/readiness probe failures

3. **Review pod logs for errors:**

   ```bash
   oc logs <POD_NAME> -n openshift-logging
   ```

   Look for startup errors, configuration issues, or crash messages.

4. **Verify the metrics endpoint is accessible:**

   If the pod is running, test if the metrics endpoint responds:

   ```bash
   oc exec <POD_NAME> -n openshift-logging -- curl -s http://localhost:8686/metrics
   ```

   The collector typically exposes metrics on port 8686. If this command fails
   or times out, the metrics server within the collector is not functioning.

5. **Check network policies:**

   Verify that network policies allow Prometheus to scrape the collector:

   ```bash
   oc get networkpolicies -n openshift-logging
   ```

   Ensure policies permit ingress traffic from the monitoring namespace (typically
   `openshift-monitoring`) to the collector pods on the metrics port.

6. **Check ServiceMonitor configuration:**

   Verify the ServiceMonitor for the collector exists and is correctly configured:

   ```bash
   oc get servicemonitor -n openshift-logging
   oc describe servicemonitor <SERVICEMONITOR_NAME> -n openshift-logging
   ```

   Ensure the ServiceMonitor's selector matches the collector service labels.

7. **Check the collector service:**

   Verify the service exists and has endpoints:

   ```bash
   oc get service -n openshift-logging -l app.kubernetes.io/component=collector
   oc get endpoints -n openshift-logging -l app.kubernetes.io/component=collector
   ```

   If the endpoints list is empty, the service selector may not match the pod labels.

## Mitigation

The following examples assume the collector is deployed to the 'openshift-logging'
namespace

### Restart Unhealthy Pods

If the collector pod is in a failed state (`CrashLoopBackOff`, `Error`), delete
it to trigger a restart:

```bash
oc delete pod <POD_NAME> -n openshift-logging
```

The DaemonSet or Deployment will automatically recreate the pod.

### Check Resource Constraints

If pod events indicate resource issues (e.g., "Pod OOMKilled"), check resource
limits and node capacity:

1. Review the pod's resource requests and limits:

   ```bash
   oc get pod <POD_NAME> -n openshift-logging -o jsonpath='{.spec.containers[*].resources}'
   ```

2. Check node resource availability:

   ```bash
   oc describe node <NODE_NAME>
   ```

3. If necessary, adjust the collector's resource limits in the
   `ClusterLogForwarder` custom resource.

### Verify ClusterLogForwarder Configuration

If the collector is failing to start due to configuration errors:

1. Review the `ClusterLogForwarder` custom resource:

   ```bash
   oc get clusterlogforwarder -A
   oc describe clusterlogforwarder <NAME> -n <NAMESPACE>
   ```

2. Check the generated collector configuration for errors:

   ```bash
   oc get configmap <COLLECTOR_CONFIG_MAP> -n openshift-logging -o yaml
   ```

3. Correct any configuration issues and the collector should restart automatically.

NOTE: Manual corrections to the configuration should be reported as an issue to
Red Hat against the cluster-logging-operator

### Check for Node-Level Issues

If the collector pod cannot be scheduled or is stuck in `Pending`:

1. Verify the node is ready:

   ```bash
   oc get nodes
   ```

2. Check for node taints that might prevent scheduling:

   ```bash
   oc describe node <NODE_NAME> | grep -i taint
   ```

3. Ensure the node has sufficient resources (CPU, memory, disk) to run the
   collector pod.

### Network Connectivity Issues

If the pod is running but Prometheus cannot reach it:

1. Verify the metrics port is exposed in the pod specification.
2. Check for network policies blocking traffic from the `openshift-monitoring`
   namespace.
3. Ensure the Service and ServiceMonitor configurations are correct.

## Notes

- The collector is typically deployed as a DaemonSet, so there should be one
  collector pod per node in the cluster.
- The default metrics port for the collector is 8686.
- This alert specifically monitors the `up` metric with labels
  `app_kubernetes_io_component = "collector"` and
  `app_kubernetes_io_part_of = "cluster-logging"`.
