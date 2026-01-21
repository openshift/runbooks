# NetObservLokiError

## Meaning

The `NetObservLokiError` alert is an operational alert that triggers when
Network Observability's flowlogs-pipeline is unable to write flow data to Loki
and is consequently dropping flows. This indicates that the Loki storage
backend is experiencing issues or is unreachable.

Unlike other Network Observability alerts, `NetObservLokiError` does not have
multiple variants based on grouping or severity levels. It is a single
critical alert that indicates a failure in the flow storage pipeline.

When flows cannot be written to Loki, they are dropped to prevent memory
exhaustion in the flowlogs-pipeline. This results in:

- Incomplete network traffic data in the Network Observability Console
- Gaps in historical network flow records
- Inability to investigate past network flows

**Note:** This is an operational alert that monitors the health of Network
Observability's storage integration, not the health of cluster network
traffic. This alert only applies when using Loki as the flow storage backend.

### Configuration limitations

Like other Network Observability operational alerts, `NetObservLokiError`
cannot be configured, other than being disabled:

- It cannot be converted to recording mode - it is always an alert
- It does not support thresholds - it fires when Loki write errors occur
  (> 0 drops) consistently after 10 minutes
- It does not support grouping - it is a global cluster-wide operational
  alert
- It cannot have variants - there is only one alert instance

The alert triggers with this hardcoded PromQL expression:

```promql
sum(rate(netobserv_loki_dropped_entries_total[1m])) > 0
```

This design is intentional because Loki write errors indicate **data loss**
and should always generate alerts rather than being silently tracked as
metrics.

### Disable this alert entirely

We do not recommend disabling this alert as it indicates data loss.
However, if needed:

```bash
oc edit flowcollector cluster
```

Add NetObservLokiError to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - NetObservLokiError
```

For more information on Network Observability Loki integration and
troubleshooting, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/troubleshooting-network-observability.html).

## Impact

When flows are being dropped due to Loki errors, the impact is:

- Loss of network flow data (gaps in historical records)
- Incomplete network traffic visualization in the Console
- Missing network metrics that depend on flow data
- Inability to audit or investigate network events during the error period
- Potential compliance violations if flow logging is required
- Reduced effectiveness of network-based alerts (may miss events)

Unlike `NetObservNoFlows`, the eBPF agents are still collecting flows and some
data may still be visible (e.g., metrics exported directly to Prometheus), but
detailed flow records are being lost.

## Diagnosis

Failing to store flows in Loki might be due to Loki being unavailable or
wrongly configured.

Make sure that Loki is up and running. If you installed it with the Loki
Operator, make sure that your `LokiStack` resource is ready.

Further analysis can be done by checking the logs in `flowlogs-pipeline`:

```bash
oc get pods -n netobserv -l app=flowlogs-pipeline
oc logs -n netobserv <POD>
```

Check for any existing network policies, both in `flowlogs-pipeline`
namespace and in the Loki namespace, if different, that could be blocking
the traffic.

## Mitigation

If the `flowslogs-pipeline` logs show connectivity issues despite Loki
being up and running, make sure that `FlowCollector` is configured
accordingly to your Loki installation. The configuration related to Loki
is defined under `spec.loki`.

If the logs show errors related to rate limiting, this is more likely an
issue related to `LokiStack` sizing.

Refer to the documentation on Loki configuration:
- [Specifically for Network
  Observability](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/installing-network-observability-operators#network-observability-loki-installation_network_observability)
- [Loki Operator general
  documentation](https://docs.redhat.com/en/documentation/red_hat_openshift_logging/latest/html/configuring_logging/configuring-lokistack-storage)
  (note that this documentation is primarily intended for OpenShift
  Logging, but largely applies to Network Observability as well)
- [Troubleshooting rate-limit error](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/installing-troubleshooting#network-observability-troubleshooting-loki-tenant-rate-limit_network-observability-troubleshooting)

If you find that a network policy is blocking the traffic:
Network policies are automatically installed when
`spec.networkPolicy.enable` is `true` in `FlowCollector`. It should not
block the traffic to Loki, however if you find that it is, you can either
configure it with additional allowed namespaces (via
`spec.networkPolicy.additionalNamespaces`), or disable it entirely and
write your own policy instead. We would also kindly ask you to report the
issue to the maintainers team.

If you have not installed Loki and don't intend to do so, you should
disable it in `FlowCollector`:

```bash
oc edit flowcollector cluster
```

Disable Loki:

```yaml
spec:
  loki:
    enable: false
```

For other mitigation strategies and solutions, refer to the [Troubleshooting Network
Observability](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/installing-troubleshooting)
documentation.
