# NetpolDenied

## Meaning

The `NetpolDenied` alert template is triggered when Network Observability
detects a high percentage of traffic being denied by Kubernetes
NetworkPolicies. This template can generate multiple alert instances depending
on how it's configured in the FlowCollector custom resource.

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

The alert fires when the percentage of connections denied by NetworkPolicies
exceeds the configured threshold. These denials are detected via eBPF hooks in
the OVN-Kubernetes network plugin.

**Note:** This alert requires the `NetworkEvents` agent feature to be enabled
in the FlowCollector configuration.

### Switch to metric-only mode (alternative to alerts)

If you want to monitor NetworkPolicy denials in the Network Health dashboard
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
      - template: NetpolDenied
        mode: MetricOnly
        variants:
        - groupBy: Namespace
          thresholds:
            info: "10"
            warning: "25"
            critical: "50"
```

In metric-only mode:

- NetworkPolicy denial violations remain visible in the **Network Health**
  dashboard
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
      - template: NetpolDenied
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "20"
            warning: "40"
            critical: "70"
```

### Disable this alert entirely

To completely disable NetpolDenied alerts:

```console
$ oc edit flowcollector cluster
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
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/observing-network-traffic.html).

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

For detailed diagnosis steps, refer to:

- [Network policy](https://docs.redhat.com/en/documentation/openshift_container_platform/4.15/html/networking/network-policy)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/4.18#Networking)

## Mitigation

For mitigation strategies and solutions, refer to:

- [Network policy](https://docs.redhat.com/en/documentation/openshift_container_platform/4.15/html/networking/network-policy)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/4.18#Networking)
