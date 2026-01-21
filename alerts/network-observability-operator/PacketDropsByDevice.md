# PacketDropsByDevice

## Meaning

The `PacketDropsByDevice` alert template is triggered when Network
Observability detects a high percentage of packet drops at the device level
based on statistics from `/proc/net/dev`. This template can generate multiple
alert instances depending on how it's configured in the FlowCollector custom
resource.

**Possible alert variants:**

- `PacketDropsByDevice_Critical` - Global cluster-wide device packet drop
  rate exceeds critical threshold
- `PacketDropsByDevice_Warning` - Global cluster-wide device packet drop rate
  exceeds warning threshold
- `PacketDropsByDevice_Info` - Global cluster-wide device packet drop rate
  exceeds info threshold
- `PacketDropsByDevice_PerNode{Critical,Warning,Info}` - Device packet drop
  rate on a specific node exceeds threshold

Unlike `PacketDropsByKernel` which tracks drops in the kernel network stack
via eBPF, `PacketDropsByDevice` monitors drops reported by network interfaces
themselves in `/proc/net/dev`. These are typically hardware-level or
driver-level drops.

**Note:** This alert does NOT require the `PacketDrop` agent feature. It uses
standard Linux network interface statistics available on all nodes.

### Switch to recording mode (alternative to alerts)

If you want to monitor device packet drops in the Network Health dashboard
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
      - template: PacketDropsByDevice
        mode: Recording
        variants:
        - groupBy: Node
          thresholds:
            info: "2"
            warning: "5"
```

In recording mode:

- Packet drop violations remain visible in the **Network Health** dashboard
- No Prometheus alerts are generated
- Metrics are still calculated and stored as recording rules
- Useful for teams that prefer passive monitoring without alert fatigue

This is particularly helpful for SRE teams who want visibility into network
health without being overwhelmed by alerts for every threshold violation.

### Adjust alert thresholds

If the alert is firing too frequently due to low thresholds, you can adjust
them:

```bash
oc edit flowcollector cluster
```

Modify the `spec.processor.metrics.healthRules` section:

```yaml
spec:
  processor:
    metrics:
      healthRules:
      - template: PacketDropsByDevice
        mode: Alert
        variants:
        - groupBy: Node
          thresholds:
            info: "5"
            warning: "10"
```

### Disable this alert entirely

To completely disable PacketDropsByDevice alerts:

```bash
oc edit flowcollector cluster
```

Add PacketDropsByDevice to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - PacketDropsByDevice
```

For more information on configuring Network Observability alerts, see the
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

## Impact

Device-level packet drops can indicate:

- Physical network issues (bad cables, faulty NICs, switch problems)
- Network interface driver issues
- Hardware buffer overruns
- Link saturation or bandwidth exhaustion
- Incompatible network settings (MTU mismatches, flow control)

High device drop rates can cause:

- Severe performance degradation
- Unstable connections
- Application timeouts and failures
- Degraded cluster control plane communication
- Storage or database performance issues if using network storage

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

- [Reducing packet drops in OVS](https://access.redhat.com/solutions/5666711)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
