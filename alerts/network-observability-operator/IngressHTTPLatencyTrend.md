# IngressHTTPLatencyTrend

## Meaning

The `IngressHTTPLatencyTrend` health rule template is triggered when Network
Observability detects a significant increase in HAProxy ingress response
latency compared to a baseline from the past. This is a trend-based alert that
compares current latency against historical values.

By default, the baseline is calculated from metrics taken 1 day ago (configurable
via `trendOffset`) averaged over a 1-hour window (configurable via `trendDuration`).

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule requires HAProxy to export the
`haproxy_server_http_average_response_latency_milliseconds` metric. This
metric is available in all recent versions of OpenShift's HAProxy router.

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
      - template: IngressHTTPLatencyTrend
        mode: Recording
        variants:
        - thresholds:
            info: "100"
            warning: "200"
          groupBy: Namespace
          trendOffset: 24h
          trendDuration: 1h
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
