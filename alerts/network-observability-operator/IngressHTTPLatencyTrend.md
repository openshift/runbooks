# IngressHTTPLatencyTrend

## Meaning

The `IngressHTTPLatencyTrend` alert template is triggered when Network
Observability detects a significant increase in HAProxy ingress response
latency compared to a baseline from the past. This is a trend-based alert that
compares current latency against historical values. This template can generate
multiple alert instances depending on how it's configured in the FlowCollector
custom resource.

**Possible alert variants:**

- `IngressHTTPLatencyTrend_Critical` - Global cluster-wide ingress latency
  increase exceeds critical threshold compared to baseline
- `IngressHTTPLatencyTrend_Warning` - Global cluster-wide ingress latency increase
  exceeds warning threshold compared to baseline
- `IngressHTTPLatencyTrend_Info` - Global cluster-wide ingress latency increase
  exceeds info threshold compared to baseline
- `IngressHTTPLatencyTrend_PerDstNamespace{Critical,Warning,Info}` - Ingress
  latency increase for traffic destined to a specific namespace exceeds
  threshold

The alert fires when the current HAProxy ingress latency exceeds the baseline
by the configured percentage threshold. By default, the baseline is calculated
from metrics taken 1 day ago (configurable via `trendOffset`) averaged over a
1-hour window (configurable via `trendDuration`). Default thresholds are 100%
for info level and 200% for warning level.

**Note:** This alert requires HAProxy to export the
`haproxy_server_http_average_response_latency_milliseconds` metric. This
metric is available in recent versions of OpenShift's HAProxy router.

### Switch to recording mode (alternative to alerts)

If you want to monitor ingress latency trends in the Network Health dashboard
without generating Prometheus alerts, you can change the health rule to
recording mode:

```bash
oc edit flowcollector cluster
```

Change the mode from `Alert` to `Recording`:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: IngressHTTPLatencyTrend
        mode: Recording
        variants:
        - groupBy: Namespace
          thresholds:
            info: "100"
            warning: "200"
            critical: "300"
          trendOffset: 24h
          trendDuration: 1h
```

In recording mode:

- Ingress latency trend violations remain visible in the **Network Health**
  dashboard
- No Prometheus alerts are generated
- Metrics are still calculated and stored as recording rules
- Useful for teams that prefer passive monitoring without alert fatigue

This is particularly helpful for SRE teams who want visibility into network
health without being overwhelmed by alerts for every threshold violation.

### Adjust alert thresholds

The thresholds are expressed in percentage of increase. If the alert is firing
too frequently due to low thresholds, you can adjust them:

```bash
oc edit flowcollector cluster
```

Modify the `spec.processor.metrics.healthRules` section:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: IngressHTTPLatencyTrend
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "150"     
            warning: "300"  
            critical: "500"
          trendOffset: 24h
          trendDuration: 1h
```

### Adjust baseline parameters

You can also adjust the baseline calculation parameters to be more or less
sensitive to changes:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: IngressHTTPLatencyTrend
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "100"
            warning: "200"
          trendOffset: 1h   
          trendDuration: 30m 
```

### Disable this alert entirely

To completely disable IngressHTTPLatencyTrend alerts:

```bash
oc edit flowcollector cluster
```

Add IngressHTTPLatencyTrend to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - IngressHTTPLatencyTrend
```

For more information on configuring Network Observability alerts and
understanding network latency, see the
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

## Impact

Increasing HAProxy ingress latency trends can indicate performance degradation
affecting user-facing applications:

- **Slow application response times**: Users experience delays when accessing
  applications through Routes
- **Degraded user experience**: Page loads, API calls, and transactions take
  longer than expected
- **Backend performance issues**: Applications or databases are responding
  slowly
- **Timeout errors**: Slow responses may eventually lead to timeouts
- **Reduced throughput**: High latency can reduce overall request processing
  capacity

Since this alert tracks trends rather than absolute values, it helps detect
gradual performance degradation that might otherwise go unnoticed until it
becomes severe.

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

Possible causes and mitigations include:

- Scaling HAProxy router pods or adjusting resource limits
- Scaling backend application pods if they're overloaded
- Checking for slow database queries or external service issues
- Reviewing HAProxy router logs and metrics for performance bottlenecks
- Checking node resources and network health
- Adjusting `trendOffset` or `trendDuration` if baseline is not representative

For mitigation strategies and solutions, refer to:

- [OpenShift Performance and Scalability](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/scalability_and_performance/index)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
