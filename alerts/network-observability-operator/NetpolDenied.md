# NetpolDenied

## Meaning

The `NetpolDenied` health rule template is triggered when Network Observability
detects a high percentage of traffic being denied by Kubernetes
NetworkPolicies.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule requires the `NetworkEvents` agent feature to be enabled
in the `FlowCollector` configuration, and must be used with the OVN-Kubernetes
Observability feature.

**Possible alert variants:**

- `NetpolDenied_Critical` - Global cluster-wide NetworkPolicy denial rate
  exceeds critical threshold
- `NetpolDenied_Warning` - Global cluster-wide NetworkPolicy denial rate
  exceeds warning threshold
- `NetpolDenied_Info` - Global cluster-wide NetworkPolicy denial rate exceeds
  info threshold
- `NetpolDenied_PerDstNamespace{Critical,Warning,Info}` - NetworkPolicy
  denial rate for traffic destined to a specific namespace exceeds threshold
- `NetpolDenied_PerSrcNamespace{Critical,Warning,Info}` - NetworkPolicy
  denial rate for traffic originating from a specific namespace exceeds
  threshold
- `NetpolDenied_PerDstNode{Critical,Warning,Info}` - NetworkPolicy denial
  rate for traffic destined to a specific node exceeds threshold
- `NetpolDenied_PerSrcNode{Critical,Warning,Info}` - NetworkPolicy denial
  rate for traffic originating from a specific node exceeds threshold
- `NetpolDenied_PerDstWorkload{Critical,Warning,Info}` - NetworkPolicy denial
  rate for traffic destined to a specific workload exceeds threshold
- `NetpolDenied_PerSrcWorkload{Critical,Warning,Info}` - NetworkPolicy denial
  rate for traffic originating from a specific workload exceeds threshold

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
      - template: NetpolDenied
        mode: Recording
        variants:
        - thresholds:
            info: "5"
            warning: "10"
          groupBy: Namespace
```

### Disable this alert entirely

To completely disable NetpolDenied alerts:

```bash
oc edit flowcollector cluster
```

Add NetpolDenied to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - NetpolDenied
```

For more information on NetworkPolicies and Network Observability, see the
[NetworkPolicy documentation](https://docs.openshift.com/container-platform/latest/networking/network_policy/about-network-policy.html)
and
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

## Impact

NetworkPolicy denials can indicate:

- Misconfigured NetworkPolicies blocking legitimate traffic
- Applications unable to communicate with required services
- Security policies working as intended (in which case the alert may be
  noise)
- Changes to NetworkPolicies causing unintended blocking
- Pods attempting unauthorized access (security event)

High denial rates can lead to:

- Application failures and service disruptions
- Failed health checks causing pod restarts
- Broken microservice communication
- Inability to access external services or APIs
- Degraded cluster functionality if control plane traffic is affected

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

- [Network policy](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_security/network-policy)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
