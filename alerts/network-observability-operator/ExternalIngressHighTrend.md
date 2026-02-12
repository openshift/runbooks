# ExternalIngressHighTrend

## Meaning

The `ExternalIngressHighTrend` health rule template is triggered when Network
Observability detects a significant increase in inbound traffic from external
networks (outside the cluster) compared to a baseline from the past. This is a
trend-based alert that compares current ingress traffic against historical
values. By default, the baseline is calculated from metrics taken 1 day ago
(configurable via `trendOffset`) averaged over a 1-hour window (configurable
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

- `ExternalIngressHighTrend_Critical` - Global cluster-wide external ingress
  traffic increase exceeds critical threshold compared to baseline
- `ExternalIngressHighTrend_Warning` - Global cluster-wide external ingress
  traffic increase exceeds warning threshold compared to baseline
- `ExternalIngressHighTrend_Info` - Global cluster-wide external ingress
  traffic increase exceeds info threshold compared to baseline
- `ExternalIngressHighTrend_PerDstNamespace{Critical,Warning,Info}` - External
  ingress traffic increase to a specific namespace exceeds threshold
- `ExternalIngressHighTrend_PerDstNode{Critical,Warning,Info}` - External
  ingress traffic increase to a specific node exceeds threshold
- `ExternalIngressHighTrend_PerDstWorkload{Critical,Warning,Info}` - External
  ingress traffic increase to a specific workload exceeds threshold

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
      - template: ExternalIngressHighTrend
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

To completely disable ExternalIngressHighTrend alerts:

```bash
oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic)
and
[Ingress Operator documentation](https://docs.openshift.com/container-platform/latest/networking/ingress-operator.html).

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
  Policies](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_security/network-policy)
  or [Admin Network
  Policies](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_security/admin-network-policy).
- Configuring [User-Defined
  Networks](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/multiple_networks/primary-networks)
  for network segmentation.
- Configuring your [Ingress Cluster
  traffic](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/ingress_and_load_balancing/configuring-ingress-cluster-traffic).
- Using [Red Hat OpenShift Service
  Mesh](https://docs.redhat.com/en/documentation/red_hat_openshift_service_mesh/latest),
  for instance to configure rate-limits on ingress.

For a comprehensive documentation, refer to the
[OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
documentation.
