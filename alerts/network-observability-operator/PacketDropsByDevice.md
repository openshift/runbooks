# PacketDropsByDevice

## Meaning

The `PacketDropsByDevice` health rule template is triggered when Network
Observability detects a high percentage of packet drops at the device level
based on statistics from `/proc/net/dev`. These are typically hardware-level or
driver-level drops.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule does NOT require the `PacketDrop` agent feature. It uses
standard Linux network interface statistics available on all nodes.

**Possible alert variants:**

- `PacketDropsByDevice_Critical` - Global cluster-wide device packet drop
  rate exceeds critical threshold
- `PacketDropsByDevice_Warning` - Global cluster-wide device packet drop rate
  exceeds warning threshold
- `PacketDropsByDevice_Info` - Global cluster-wide device packet drop rate
  exceeds info threshold
- `PacketDropsByDevice_PerNode{Critical,Warning,Info}` - Device packet drop
  rate on a specific node exceeds threshold

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
      - template: PacketDropsByDevice
        mode: Alert
        variants:
        - thresholds:
            info: "5"
            warning: "10"
          groupBy: Node
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
