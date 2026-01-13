# ExternalEgressHighTrend

## Meaning

The `ExternalEgressHighTrend` alert template is triggered when Network
Observability detects a significant increase in outbound traffic to external
networks (outside the cluster) compared to a baseline from the past. This is a
trend-based alert that compares current egress traffic against historical
values. This template can generate multiple alert instances depending on how
it's configured in the FlowCollector custom resource.

**Possible alert variants:**

- `ExternalEgressHighTrend_Critical` - Global cluster-wide external egress
  traffic increase exceeds critical threshold compared to baseline
- `ExternalEgressHighTrend_Warning` - Global cluster-wide external egress
  traffic increase exceeds warning threshold compared to baseline
- `ExternalEgressHighTrend_Info` - Global cluster-wide external egress traffic
  increase exceeds info threshold compared to baseline
- `ExternalEgressHighTrend_PerSrcNamespace{Critical,Warning,Info}` - External
  egress traffic increase from a specific namespace exceeds threshold
- `ExternalEgressHighTrend_PerSrcNode{Critical,Warning,Info}` - External
  egress traffic increase from a specific node exceeds threshold
- `ExternalEgressHighTrend_PerSrcWorkload{Critical,Warning,Info}` - External
  egress traffic increase from a specific workload exceeds threshold

The alert fires when the current external egress traffic volume exceeds the
baseline by the configured percentage threshold. By default, the baseline is
calculated from metrics taken 1 day ago (configurable via `trendOffset`)
averaged over a 2-hour window (configurable via `trendDuration`).

External traffic is defined as traffic where the destination IP is not within
the cluster's pod or service CIDR ranges.

**Note:** This alert does not require any specific agent feature to be
enabled.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor external egress trends in the Network Health dashboard
without generating Prometheus alerts, you can change the health rule to
metric-only mode:

```console
$ oc edit flowcollector cluster
```

Change the mode from `Alert` to `MetricOnly`:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: ExternalEgressHighTrend
        mode: MetricOnly
        variants:
        - groupBy: Namespace
          thresholds:
            info: "50"
            warning: "150"
            critical: "300"
          trendOffset: 24h
          trendDuration: 2h
```

In metric-only mode:

- External egress trend violations remain visible in the **Network Health**
  dashboard
- No Prometheus alerts are generated
- Metrics are still calculated and stored as recording rules
- Useful for teams that prefer passive monitoring without alert fatigue

This is particularly helpful for SRE teams who want visibility into network
health without being overwhelmed by alerts for every threshold violation.

### Adjust alert thresholds

If the alert is firing too frequently due to low thresholds, you can adjust
them:

```console
$ oc edit flowcollector cluster
```

Modify the `spec.processor.metrics.healthRules` section:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: ExternalEgressHighTrend
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "100"
            warning: "200"
            critical: "400"
          trendOffset: 24h
          trendDuration: 2h
```

### Disable this alert entirely

To completely disable ExternalEgressHighTrend alerts:

```console
$ oc edit flowcollector cluster
```

Add ExternalEgressHighTrend to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - ExternalEgressHighTrend
```

For more information on configuring Network Observability alerts and
monitoring external traffic, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html).

## Impact

Unexpected increases in external egress traffic can indicate:

- Application changes causing increased API calls to external services
- Data exfiltration (security concern)
- Misconfigured applications polling external services too frequently
- Workload scaling causing proportional increase in external traffic
- Backup or sync operations to external storage
- Distributed Denial of Service (DDoS) attacks originating from the cluster
- Compromised pods acting as part of a botnet

High external egress traffic can lead to:

- Increased cloud egress costs (especially in public clouds)
- Network bandwidth exhaustion
- Degraded performance for other services
- Potential security incidents if unauthorized
- Rate limiting or blocking by external APIs
- Compliance violations if sensitive data is being exfiltrated

## Diagnosis

For detailed diagnosis steps, refer to the [Troubleshooting Network
Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
documentation.

## Mitigation

For mitigation strategies and solutions, refer to the [Troubleshooting Network
Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
documentation.
