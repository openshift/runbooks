# PacketDropsByDevice

## Meaning

The `PacketDropsByDevice` alert template is triggered when Network Observability
detects a high percentage of packet drops at the device level based on
statistics from `/proc/net/dev`. This template can generate multiple alert
instances depending on how it's configured in the FlowCollector custom resource.

**Possible alert variants:**

- `PacketDropsByDevice_Critical` - Global cluster-wide device packet drop rate
  exceeds critical threshold (no grouping)
- `PacketDropsByDevice_Warning` - Global cluster-wide device packet drop rate
  exceeds warning threshold (no grouping)
- `PacketDropsByDevice_Info` - Global cluster-wide device packet drop rate
  exceeds info threshold (no grouping)
- `PacketDropsByDevice_PerNode{Critical,Warning,Info}` - Device packet drop
  rate on a specific node exceeds threshold

Unlike `PacketDropsByKernel` which tracks drops in the kernel network stack via
eBPF, `PacketDropsByDevice` monitors drops reported by network interfaces
themselves in `/proc/net/dev`. These are typically hardware-level or
driver-level drops.

**Note:** This alert does NOT require the `PacketDrop` agent feature. It uses
standard Linux network interface statistics available on all nodes.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor device packet drops in the Network Health dashboard without
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
      - template: PacketDropsByDevice
        mode: MetricOnly
        variants:
        - groupBy: Node
          thresholds:
            info: "2"
            warning: "5"
```

In metric-only mode:

- Packet drop violations remain visible in the **Network Health** dashboard
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
      - template: PacketDropsByDevice
        mode: Alert
        variants:
        - groupBy: Node
          thresholds:
            info: "5"       # Increased from 2
            warning: "10"   # Increased from 5
```

### Disable this alert entirely

To completely disable PacketDropsByDevice alerts:

```console
$ oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html).

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

TBD

## Mitigation

TBD
