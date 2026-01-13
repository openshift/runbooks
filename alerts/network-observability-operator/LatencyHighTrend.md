# LatencyHighTrend

## Meaning

The `LatencyHighTrend` alert template is triggered when Network Observability
detects a significant increase in TCP latency (Round-Trip Time) compared to a
baseline from the past. This is a trend-based alert that compares current
latency against historical values. This template can generate multiple alert
instances depending on how it's configured in the FlowCollector custom
resource.

**Possible alert variants:**

- `LatencyHighTrend_Critical` - Global cluster-wide latency increase exceeds
  critical threshold compared to baseline
- `LatencyHighTrend_Warning` - Global cluster-wide latency increase exceeds
  warning threshold compared to baseline
- `LatencyHighTrend_Info` - Global cluster-wide latency increase exceeds info
  threshold compared to baseline
- `LatencyHighTrend_PerDstNamespace{Critical,Warning,Info}` - Latency
  increase for traffic destined to a specific namespace exceeds threshold
- `LatencyHighTrend_PerSrcNamespace{Critical,Warning,Info}` - Latency
  increase for traffic originating from a specific namespace exceeds threshold
- `LatencyHighTrend_PerDstNode{Critical,Warning,Info}` - Latency increase for
  traffic destined to a specific node exceeds threshold
- `LatencyHighTrend_PerSrcNode{Critical,Warning,Info}` - Latency increase for
  traffic originating from a specific node exceeds threshold
- `LatencyHighTrend_PerDstWorkload{Critical,Warning,Info}` - Latency increase
  for traffic destined to a specific workload exceeds threshold
- `LatencyHighTrend_PerSrcWorkload{Critical,Warning,Info}` - Latency increase
  for traffic originating from a specific workload exceeds threshold

The alert fires when the current latency exceeds the baseline by the
configured percentage threshold. By default, the baseline is calculated from
metrics taken 1 day ago (configurable via `trendOffset`) averaged over a
2-hour window (configurable via `trendDuration`).

**Note:** This alert requires the `FlowRTT` agent feature to be enabled in the
FlowCollector configuration.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor latency trends in the Network Health dashboard without
generating Prometheus alerts, you can change the health rule to metric-only
mode:

```console
$ oc edit flowcollector cluster
```

Change the mode from `Alert` to `MetricOnly`:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: LatencyHighTrend
        mode: MetricOnly
        variants:
        - groupBy: Namespace
          thresholds:
            info: "30"
            warning: "75"
            critical: "150"
          trendOffset: 24h
          trendDuration: 2h
```

In metric-only mode:

- Latency trend violations remain visible in the **Network Health** dashboard
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
      - template: LatencyHighTrend
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "50"
            warning: "100"
            critical: "200"
          trendOffset: 24h
          trendDuration: 2h
```

### Disable this alert entirely

To completely disable LatencyHighTrend alerts:

```console
$ oc edit flowcollector cluster
```

Add LatencyHighTrend to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - LatencyHighTrend
```

For more information on configuring Network Observability alerts and
understanding network latency, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html).

## Impact

Increasing TCP latency trends can indicate:

- Network congestion or saturation
- Routing issues or suboptimal paths
- Hardware problems (failing NICs, switches)
- CPU or resource contention on nodes
- External network degradation
- Application performance degradation

High or increasing latency can lead to:

- Slow application response times
- Degraded user experience
- Timeout errors in distributed systems
- Reduced throughput (TCP congestion control reduces transmission rate)
- Failed health checks causing unnecessary pod restarts
- Database query slowdowns
- Inefficient microservice communication

## Diagnosis

For detailed diagnosis steps, refer to:

- [Troubleshooting latency in OpenShift](https://access.redhat.com/solutions/4255351)
- [Troubleshooting Network Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)

## Mitigation

For mitigation strategies and solutions, refer to:

- [Troubleshooting latency in OpenShift](https://access.redhat.com/solutions/4255351)
- [Troubleshooting Network Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
