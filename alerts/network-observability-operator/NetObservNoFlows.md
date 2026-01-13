# NetObservNoFlows

## Meaning

The `NetObservNoFlows` alert is an operational alert that triggers when
Network Observability is not receiving any network flow data for a certain
period. This indicates that the flow collection pipeline has stopped working,
preventing network traffic visibility.

Unlike other Network Observability alerts, `NetObservNoFlows` does not have
multiple variants based on grouping or severity levels. It is a single
critical alert that indicates a complete failure in flow collection.

This alert monitors the `flowlogs-pipeline` component which processes and
exports network flows. When no flows are received, it means either:

- The eBPF agents on nodes are not collecting flows
- The flowlogs-pipeline is not receiving data from agents
- The flowlogs-pipeline has crashed or stopped processing

**Note:** This is an operational alert that monitors the health of Network
Observability itself, not the health of cluster network traffic.

### Configuration limitations

**Important:** `NetObservNoFlows` is an **alert-only template** that cannot be
configured as a health rule. Unlike other Network Observability alerts (such
as DNSErrors or PacketDropsByKernel), this alert:

- **Cannot be converted to metric-only mode** - it is always an alert
- **Does not support thresholds** - it fires when no flows are received (flow
  rate == 0)
- **Does not support grouping** - it is a global cluster-wide operational
  alert
- **Cannot have variants** - there is only one alert instance

The alert triggers with this hardcoded PromQL expression:
```promql
sum(rate(netobserv_ingest_flows_processed[1m])) == 0
```

This design is intentional because no flows being received indicates a
**complete failure** of Network Observability and should always generate
alerts.

### Disable this alert entirely

This alert should generally NOT be disabled as it indicates a critical
failure. However, if needed:

```console
$ oc edit flowcollector cluster
```

Add NetObservNoFlows to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - NetObservNoFlows
```

For more information on Network Observability operational health and
troubleshooting, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/troubleshooting-network-observability.html).

## Impact

When no flows are being observed, the impact is:

- Complete loss of network traffic visibility in the cluster
- Network Observability Console shows no data or stale data
- Network Health dashboard shows no metrics
- All network-based alerts will not fire (blind to network issues)
- Unable to troubleshoot network problems
- Loss of network traffic audit trail
- Compliance violations if network monitoring is required

This is a critical operational issue that requires immediate attention, as it
means Network Observability is completely non-functional.

## Diagnosis

For detailed diagnosis steps, refer to the [Troubleshooting Network
Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
documentation.

## Mitigation

For mitigation strategies and solutions, refer to the [Troubleshooting Network
Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
documentation.
