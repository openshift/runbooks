# NetObservLokiError

## Meaning

The `NetObservLokiError` alert is an operational alert that triggers when
Network Observability's flowlogs-pipeline is unable to write flow data to Loki
and is consequently dropping flows. This indicates that the Loki storage backend
is experiencing issues or is unreachable.

Unlike other Network Observability alerts, `NetObservLokiError` does not have
multiple variants based on grouping or severity levels. It is a single critical
alert that indicates a failure in the flow storage pipeline.

When flows cannot be written to Loki, they are dropped to prevent memory
exhaustion in the flowlogs-pipeline. This results in:

- Incomplete network traffic data in the Network Observability Console
- Gaps in historical network flow records
- Missing network metrics in Prometheus
- Inability to investigate past network events

**Note:** This is an operational alert that monitors the health of Network
Observability's storage integration, not the health of cluster network traffic.
This alert only applies when using Loki as the flow storage backend.

### Disable this alert entirely

This alert should generally NOT be disabled as it indicates data loss. However,
if needed:

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

TBD

## Mitigation

TBD
