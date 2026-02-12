# DNSErrors

## Meaning

The `DNSErrors` health rule template is triggered when Network Observability
detects a high percentage of DNS errors (excluding NX_DOMAIN errors, which
have their own alert template). DNS errors are tracked only in return traffic
(responses from DNS servers).

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule requires the `DNSTracking` agent feature to be enabled
in the `FlowCollector` configuration.

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
      - template: DNSErrors
        mode: Alert
        variants:
        - thresholds:
            warning: "5"
        - thresholds:
            info: "5"
            warning: "10"
          groupBy: Namespace
```

If you prefer to switch to the recording mode with `mode: Recording`:

- DNS error violations remain visible in the **Network Health** dashboard
- No Prometheus alerts are generated
- Metrics are still calculated and stored as recording rules
- Useful for teams that prefer passive monitoring without alert fatigue

### Disable this alert entirely

To completely disable DNSErrors alerts:

```bash
oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

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

- [Troubleshooting DNS in OpenShift](https://access.redhat.com/solutions/3804501)
- [OpenShift DNS Operator](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/networking_operators/dns-operator)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
