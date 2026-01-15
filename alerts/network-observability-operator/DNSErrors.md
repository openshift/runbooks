# DNSErrors

## Meaning

The `DNSErrors` alert template is triggered when Network Observability
detects a high percentage of DNS errors (excluding NX_DOMAIN errors, which
have their own alert template). This template can generate multiple alert
instances depending on how it's configured in the FlowCollector custom
resource.

**Possible alert variants:**

- `DNSErrors_Critical` - Global cluster-wide DNS error rate exceeds critical
  threshold
- `DNSErrors_Warning` - Global cluster-wide DNS error rate exceeds warning
  threshold
- `DNSErrors_Info` - Global cluster-wide DNS error rate exceeds info
  threshold
- `DNSErrors_PerDstNamespace{Critical,Warning,Info}` - DNS error rate for
  traffic destined to a specific namespace exceeds threshold
- `DNSErrors_PerSrcNamespace{Critical,Warning,Info}` - DNS error rate for
  traffic originating from a specific namespace exceeds threshold
- `DNSErrors_PerDstNode{Critical,Warning,Info}` - DNS error rate for traffic
  destined to a specific node exceeds threshold
- `DNSErrors_PerSrcNode{Critical,Warning,Info}` - DNS error rate for traffic
  originating from a specific node exceeds threshold
- `DNSErrors_PerDstWorkload{Critical,Warning,Info}` - DNS error rate for
  traffic destined to a specific workload exceeds threshold
- `DNSErrors_PerSrcWorkload{Critical,Warning,Info}` - DNS error rate for
  traffic originating from a specific workload exceeds threshold

The alert fires when the percentage of DNS errors exceeds the configured
threshold. DNS errors are tracked only in return traffic (responses from DNS
servers).

**Note:** This alert requires the `DNSTracking` agent feature to be enabled
in the FlowCollector configuration.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor DNS errors in the Network Health dashboard without
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
      - template: DNSErrors
        mode: MetricOnly
        variants:
        - groupBy: Namespace
          thresholds:
            info: "5"
            warning: "20"
            critical: "50"
```

In metric-only mode:

- DNS error violations remain visible in the **Network Health** dashboard
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
      - template: DNSErrors
        mode: alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "10"      # Increased from 5
            warning: "25"   # Increased from 20
            critical: "60"  # Increased from 50
```

### Disable this alert entirely

To completely disable DNSErrors alerts:

```console
$ oc edit flowcollector cluster
```

Add DNSErrors to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - DNSErrors
```

For more information on configuring Network Observability alerts, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html).

## Impact

DNS errors prevent applications from resolving domain names to IP addresses,
which can lead to:

- Application failures and connection timeouts
- Inability to reach external services or internal cluster services
- Degraded user experience due to failed service lookups
- Cascading failures in microservices architectures that depend on service
  discovery

Even a moderate percentage of DNS errors can significantly impact application
reliability, as DNS resolution is typically on the critical path for most
network operations.

## Diagnosis

When this alert fires, you can investigate further by using the Network Observability interface:

1. **Navigate to alert details**: Click on the alert in the Network Health dashboard to view specific details of the alert.

2. **Navigate to network traffic**: From the alert details, you can navigate to the Network Traffic view to examine the specific flows that are related to this alert. This allows you to see:
   - Source and destination of the traffic
   - Detailed flow information

For additional troubleshooting resources, refer to the documentation links in the Mitigation section below.

## Mitigation

For mitigation strategies and solutions, refer to:

- [Troubleshooting DNS in OpenShift](https://access.redhat.com/solutions/3804501)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/4.18#Networking)
