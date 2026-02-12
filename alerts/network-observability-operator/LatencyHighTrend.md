# LatencyHighTrend

## Meaning

The `LatencyHighTrend` health rule template is triggered when Network Observability
detects a significant increase in TCP latency (Round-Trip Time) compared to a
baseline from the past. This is a trend-based alert that compares current
latency against historical values.

By default, the baseline is calculated from metrics taken 1 day ago (configurable
via `trendOffset`) averaged over a 1-hour window (configurable via `trendDuration`).

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule requires the `FlowRTT` agent feature to be enabled in the
FlowCollector configuration.

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

### Default definition

You can override the default definition by editing the `FlowCollector` resource:

```bash
oc edit flowcollector cluster
```

Insert these default values, and edit them as desired:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: LatencyHighTrend
        mode: Recording
        variants:
        - thresholds:
            info: "100"
          groupBy: Namespace
          trendOffset: 24h
          trendDuration: 1h
```

### Disable this alert entirely

To completely disable LatencyHighTrend alerts:

```bash
oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

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

For mitigation strategies and solutions, refer to:

- [Troubleshooting latency in OpenShift](https://access.redhat.com/solutions/4255351)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
