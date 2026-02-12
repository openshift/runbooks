# ExternalEgressHighTrend

## Meaning

The `ExternalEgressHighTrend` health rule template is triggered when Network
Observability detects a significant increase in outbound traffic to external
networks (outside the cluster) compared to a baseline from the past. This is a
trend-based alert that compares current egress traffic against historical
values. By default, the baseline is calculated from metrics taken 1 day ago
(configurable via `trendOffset`) averaged over a 2-hour window (configurable
via `trendDuration`).

External traffic is defined as traffic where the source IP is not within the
cluster's node, pod or service CIDR ranges. Those ranges are defined in the
`FlowCollector` resource, under `spec.processor.subnetLabels`.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule does not require any specific agent feature to be
enabled.

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
      - template: ExternalEgressHighTrend
        mode: Recording
        variants:
        - thresholds:
            warning: "200"
          groupBy: Node
          trendOffset: 24h
          trendDuration: 1h
        - thresholds:
            info: "100"
            warning: "500"
          groupBy: Namespace
          trendOffset: 24h
          trendDuration: 1h
```

### Disable this alert entirely

To completely disable ExternalEgressHighTrend alerts:

```bash
oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

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

Depending on the cause of the traffic, the mitigation can be:

- Defining [Network
  Policies](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_security/network-policy),
  [Admin Network
  Policies](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_security/admin-network-policy)
  or configuring an [Egress
  Firewall](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_security/egress-firewall).
- Configuring [User-Defined
  Networks](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/multiple_networks/primary-networks)
  for network segmentation.
- Using [Red Hat OpenShift Service
  Mesh](https://docs.redhat.com/en/documentation/red_hat_openshift_service_mesh/latest),
  for instance to configure rate-limits on egress.

For a comprehensive documentation, refer to the
[OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
documentation.
