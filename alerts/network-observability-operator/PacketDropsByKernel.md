# PacketDropsByKernel

## Meaning

The `PacketDropsByKernel` health rule template is triggered when Network
Observability detects a high percentage of packet dropped by the Linux kernel,
based on the configured threshold. These drops are detected via eBPF kfree_skb
tracepoint which captures packets dropped in the kernel network stack.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule requires the `PacketDrop` agent feature to be enabled in
the `FlowCollector` configuration.

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
      - template: PacketDropsByKernel
        mode: Recording
        variants:
        - lowVolumeThreshold: "5"
          thresholds:
            info: "10"
            warning: "20"
          groupBy: Namespace
        - thresholds:
            info: "5"
            warning: "15"
          groupBy: Node
```

### Disable this alert entirely

To completely disable PacketDropsByKernel alerts:

```bash
oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

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

When this alert fires, you can investigate further by using the Network
Observability interface:

1. **Navigate to alert details**: Click on the alert in the Network Health
   dashboard to view specific details of the alert.

2. **Navigate to network traffic**: From the alert details, you can navigate
   to the Network Traffic view to examine the specific flows that are related
   to this alert. This allows you to see:
   - Source and destination of the traffic
   - Detailed flow information

In the Network Traffic view, to further narrow down your search when
looking for drops, you can open the _Query options_ dropdown menu, and
select flows containing drops.

The drop causes that Network Observability displays are derived directly
from the error codes provided [by the Linux
kernel](https://github.com/torvalds/linux/blob/master/include/net/dropreason-core.h#L140),
or [by
OVS](https://git.kernel.org/pub/scm/linux/kernel/git/netdev/net-next.git/tree/net/openvswitch/drop.h).
Network Observability does not provide its own interpretation of them.

When you use Network Policies, you might find the `OVS_DROP_LAST_ACTION` cause appearing
more frequently: it is set when a network policy has rejected a packet.

For additional troubleshooting resources, refer to the documentation links in
the Mitigation section below.

## Mitigation

For mitigation strategies and solutions, refer to:

- [Packet drop tracking in Network Observability](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/nw-observe-network-traffic#network-observability-pktdrop-overview_nw-observe-network-traffic)
- [Reducing packet drops in OVS](https://access.redhat.com/solutions/5666711)
- [Blog: Network Observability real-time per flow packets
  drop](https://www.redhat.com/en/blog/network-observability-real-time-per-flow-packets-drop)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
