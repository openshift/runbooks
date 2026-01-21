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

Like other Network Observability operational alerts, `NetObservLokiError` cannot be configured, other than being disabled:

- It cannot be converted to metric-only mode - it is always an alert
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

```console
$ oc edit flowcollector cluster
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

When this alert fires, you can investigate further by using the Network
Observability interface:

1. **Navigate to alert details**: Click on the alert in the Network Health
   dashboard to view specific details of the alert.

2. **Navigate to network traffic**: From the alert details, you can navigate
   to the Network Traffic view to examine the specific flows that are related
   to this alert. This allows you to see:
   - Source and destination of the traffic
   - Detailed flow information

For additional troubleshooting resources, refer to the documentation links in
the Mitigation section below.

## Mitigation

For mitigation strategies and solutions, refer to the [Troubleshooting Network
Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
documentation.
