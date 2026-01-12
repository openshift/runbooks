# ExternalIngressHighTrend

## Meaning

The `ExternalIngressHighTrend` alert template is triggered when Network
Observability detects a significant increase in inbound traffic from external
networks (outside the cluster) compared to a baseline from the past. This is a
trend-based alert that compares current ingress traffic against historical
values. This template can generate multiple alert instances depending on how
it's configured in the FlowCollector custom resource.

**Possible alert variants:**

- `ExternalIngressHighTrend_Critical` - Global cluster-wide external ingress
  traffic increase exceeds critical threshold compared to baseline (no grouping)
- `ExternalIngressHighTrend_Warning` - Global cluster-wide external ingress
  traffic increase exceeds warning threshold compared to baseline (no grouping)
- `ExternalIngressHighTrend_Info` - Global cluster-wide external ingress
  traffic increase exceeds info threshold compared to baseline (no grouping)
- `ExternalIngressHighTrend_PerDstNamespace{Critical,Warning,Info}` - External
  ingress traffic increase to a specific namespace exceeds threshold
- `ExternalIngressHighTrend_PerDstNode{Critical,Warning,Info}` - External
  ingress traffic increase to a specific node exceeds threshold
- `ExternalIngressHighTrend_PerDstWorkload{Critical,Warning,Info}` - External
  ingress traffic increase to a specific workload exceeds threshold

The alert fires when the current external ingress traffic volume exceeds the
baseline by the configured percentage threshold. By default, the baseline is
calculated from metrics taken 1 day ago (configurable via `trendOffset`)
averaged over a 2-hour window (configurable via `trendDuration`).

External traffic is defined as traffic where the source IP is not within the
cluster's pod or service CIDR ranges.

**Note:** This alert does not require any specific agent feature to be enabled.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor external ingress trends in the Network Health dashboard without
generating Prometheus alerts, you can change the health rule to metric-only mode:

```console
$ oc edit flowcollector cluster
```

Change the mode from `Alert` to `MetricOnly`:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: ExternalIngressHighTrend
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

- External ingress trend violations remain visible in the **Network Health** dashboard
- No Prometheus alerts are generated (no AlertManager notifications)
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
      - template: ExternalIngressHighTrend
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "100"      # Increased from 50
            warning: "200"   # Increased from 150
            critical: "400"  # Increased from 300
          trendOffset: 24h
          trendDuration: 2h
```

### Disable this alert entirely

To completely disable ExternalIngressHighTrend alerts:

```console
$ oc edit flowcollector cluster
```

Add ExternalIngressHighTrend to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - ExternalIngressHighTrend
```

For more information on configuring Network Observability alerts, managing
ingress traffic, and DDoS protection, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html)
and [Ingress Operator documentation](https://docs.openshift.com/container-platform/latest/networking/ingress-operator.html).

## Impact

Unexpected increases in external ingress traffic can indicate:

- Legitimate traffic surges (viral content, marketing campaigns, seasonal
  spikes)
- Distributed Denial of Service (DDoS) attacks
- Web scraping or bot traffic
- Application becoming more popular or being linked from high-traffic sites
- Load testing or stress testing
- Misconfigured external systems flooding the cluster with requests
- Retry storms from external clients

High external ingress traffic can lead to:

- Degraded application performance due to resource exhaustion
- Service outages if capacity is exceeded
- Increased cloud ingress costs (though typically lower than egress)
- Triggered autoscaling that may increase infrastructure costs
- Failed requests and poor user experience
- Database overload from increased query load
- Potential security breaches if attack traffic is not mitigated

## Diagnosis

TBD

## Mitigation

TBD
