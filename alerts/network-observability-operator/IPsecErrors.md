# IPsecErrors

## Meaning

The `IPsecErrors` alert template is triggered when Network Observability detects
a high percentage of IPsec encryption errors. This template can generate
multiple alert instances depending on how it's configured in the FlowCollector
custom resource.

**Possible alert variants:**

- `IPsecErrors_Critical` - Global cluster-wide IPsec error rate exceeds
  critical threshold (no grouping)
- `IPsecErrors_Warning` - Global cluster-wide IPsec error rate exceeds warning
  threshold (no grouping)
- `IPsecErrors_Info` - Global cluster-wide IPsec error rate exceeds info
  threshold (no grouping)
- `IPsecErrors_PerDstNamespace{Critical,Warning,Info}` - IPsec error rate for
  traffic destined to a specific namespace exceeds threshold
- `IPsecErrors_PerSrcNamespace{Critical,Warning,Info}` - IPsec error rate for
  traffic originating from a specific namespace exceeds threshold
- `IPsecErrors_PerDstNode{Critical,Warning,Info}` - IPsec error rate for
  traffic destined to a specific node exceeds threshold
- `IPsecErrors_PerSrcNode{Critical,Warning,Info}` - IPsec error rate for
  traffic originating from a specific node exceeds threshold
- `IPsecErrors_PerDstWorkload{Critical,Warning,Info}` - IPsec error rate for
  traffic destined to a specific workload exceeds threshold
- `IPsecErrors_PerSrcWorkload{Critical,Warning,Info}` - IPsec error rate for
  traffic originating from a specific workload exceeds threshold

The alert fires when the percentage of IPsec encryption/decryption errors
exceeds the configured threshold. These errors are detected via eBPF hooks in
the IPsec processing path.

**Note:** This alert requires the `IPSec` agent feature to be enabled in the
FlowCollector configuration. IPsec monitoring is typically used in clusters
with OVN-Kubernetes IPsec encryption enabled.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor IPsec errors in the Network Health dashboard without
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
      - template: IPsecErrors
        mode: MetricOnly
        variants:
        - groupBy: Node
          thresholds:
            info: "5"
            warning: "15"
            critical: "30"
```

In metric-only mode:

- IPsec error violations remain visible in the **Network Health** dashboard
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
      - template: IPsecErrors
        mode: Alert
        variants:
        - groupBy: Node
          thresholds:
            info: "10"      # Increased from 5
            warning: "25"   # Increased from 15
            critical: "50"  # Increased from 30
```

### Disable this alert entirely

To completely disable IPsecErrors alerts:

```console
$ oc edit flowcollector cluster
```

Add IPsecErrors to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - IPsecErrors
```

For more information on configuring Network Observability alerts and IPsec
encryption, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html)
and [OVN-Kubernetes IPsec documentation](https://docs.openshift.com/container-platform/latest/networking/ovn_kubernetes_network_provider/configuring-ipsec-ovn.html).

## Impact

IPsec errors indicate failures in encrypting or decrypting network traffic,
which can lead to:

- Unencrypted traffic being transmitted (security violation)
- Dropped packets that fail encryption/decryption
- Communication failures between pods or nodes
- Data integrity issues if packets are corrupted
- Compliance violations if encryption is required by policy

IPsec is critical for securing pod-to-pod and node-to-node communication in
multi-tenant or security-sensitive environments. Errors can expose sensitive
data.

## Diagnosis

TBD

## Mitigation

TBD
