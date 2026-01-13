# PacketDropsByKernel

## Meaning

The `PacketDropsByKernel` alert template is triggered when Network
Observability detects a high percentage of packet drops by the Linux kernel.
This template can generate multiple alert instances depending on how it's
configured in the FlowCollector custom resource.

**Possible alert variants:**

- `PacketDropsByKernel_Critical` - Global cluster-wide kernel packet drop
  rate exceeds critical threshold
- `PacketDropsByKernel_Warning` - Global cluster-wide kernel packet drop rate
  exceeds warning threshold
- `PacketDropsByKernel_Info` - Global cluster-wide kernel packet drop rate
  exceeds info threshold
- `PacketDropsByKernel_PerDstNamespace{Critical,Warning,Info}` - Kernel
  packet drop rate for traffic destined to a specific namespace exceeds
  threshold
- `PacketDropsByKernel_PerSrcNamespace{Critical,Warning,Info}` - Kernel
  packet drop rate for traffic originating from a specific namespace exceeds
  threshold
- `PacketDropsByKernel_PerDstNode{Critical,Warning,Info}` - Kernel packet
  drop rate for traffic destined to a specific node exceeds threshold
- `PacketDropsByKernel_PerSrcNode{Critical,Warning,Info}` - Kernel packet
  drop rate for traffic originating from a specific node exceeds threshold
- `PacketDropsByKernel_PerDstWorkload{Critical,Warning,Info}` - Kernel packet
  drop rate for traffic destined to a specific workload exceeds threshold
- `PacketDropsByKernel_PerSrcWorkload{Critical,Warning,Info}` - Kernel packet
  drop rate for traffic originating from a specific workload exceeds
  threshold

The alert fires when the percentage of packets dropped by the kernel exceeds
the configured threshold. These drops are detected via eBPF kfree_skb
tracepoint which captures packets dropped in the kernel network stack.

**Note:** This alert requires the `PacketDrop` agent feature to be enabled in
the FlowCollector configuration.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor kernel packet drops in the Network Health dashboard
without generating Prometheus alerts, you can change the health rule to
metric-only mode:

```console
$ oc edit flowcollector cluster
```

Change the mode from `Alert` to `MetricOnly`:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: PacketDropsByKernel
        mode: MetricOnly
        variants:
        - groupBy: Node
          thresholds:
            info: "5"
            warning: "10"
            critical: "20"
```

In metric-only mode:

- Packet drop violations remain visible in the **Network Health** dashboard
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
      - template: PacketDropsByKernel
        mode: Alert
        variants:
        - groupBy: Node
          thresholds:
            info: "10"
            warning: "20"
            critical: "30"
```

### Disable this alert entirely

To completely disable PacketDropsByKernel alerts:

```console
$ oc edit flowcollector cluster
```

Add PacketDropsByKernel to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - PacketDropsByKernel
```

For more information on configuring Network Observability alerts, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html).

## Impact

Packet drops by the kernel can lead to:

- Degraded application performance due to retransmissions
- Increased latency as TCP backs off and retransmits
- Reduced throughput for network-intensive applications
- Failed connections if drop rates are severe
- Poor user experience in real-time applications (VoIP, video streaming)

Even moderate packet drop rates (5-10%) can significantly impact TCP
performance due to congestion control mechanisms.

## Diagnosis

For detailed diagnosis steps, refer to:

- [Reducing packet drops in OVS](https://access.redhat.com/solutions/5666711)
- [Troubleshooting Network Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)

## Mitigation

For mitigation strategies and solutions, refer to:

- [Reducing packet drops in OVS](https://access.redhat.com/solutions/5666711)
- [Troubleshooting Network Observability](https://docs.openshift.com/container-platform/latest/observability/network_observability/troubleshooting-network-observability.html)
