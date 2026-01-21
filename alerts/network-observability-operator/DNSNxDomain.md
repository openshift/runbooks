# DNSNxDomain

## Meaning

The `DNSNxDomain` alert template is triggered when Network Observability
detects a high percentage of DNS NX_DOMAIN errors. NX_DOMAIN indicates that
the queried domain name does not exist. This template can generate multiple
alert instances depending on how it's configured in the FlowCollector custom
resource.

**Possible alert variants:**

- `DNSNxDomain_Critical` - Global cluster-wide NX_DOMAIN rate exceeds
  critical threshold
- `DNSNxDomain_Warning` - Global cluster-wide NX_DOMAIN rate exceeds warning
  threshold
- `DNSNxDomain_Info` - Global cluster-wide NX_DOMAIN rate exceeds info
  threshold
- `DNSNxDomain_PerDstNamespace{Critical,Warning,Info}` - NX_DOMAIN rate for
  traffic destined to a specific namespace exceeds threshold
- `DNSNxDomain_PerSrcNamespace{Critical,Warning,Info}` - NX_DOMAIN rate for
  traffic originating from a specific namespace exceeds threshold
- `DNSNxDomain_PerDstNode{Critical,Warning,Info}` - NX_DOMAIN rate for
  traffic destined to a specific node exceeds threshold
- `DNSNxDomain_PerSrcNode{Critical,Warning,Info}` - NX_DOMAIN rate for
  traffic originating from a specific node exceeds threshold
- `DNSNxDomain_PerDstWorkload{Critical,Warning,Info}` - NX_DOMAIN rate for
  traffic destined to a specific workload exceeds threshold
- `DNSNxDomain_PerSrcWorkload{Critical,Warning,Info}` - NX_DOMAIN rate for
  traffic originating from a specific workload exceeds threshold

NX_DOMAIN errors are tracked only in return traffic (responses from DNS
servers).

**Note:** This alert requires the `DNSTracking` agent feature to be enabled
in the FlowCollector configuration.

### Switch to recording mode (alternative to alerts)

If you want to monitor DNS NX_DOMAIN errors in the Network Health dashboard
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
      - template: DNSNxDomain
        mode: Recording
        variants:
        - groupBy: Namespace
          thresholds:
            info: "10"
            warning: "30"
            critical: "60"
```

In recording mode:

- NX_DOMAIN violations remain visible in the **Network Health** dashboard
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
      - template: DNSNxDomain
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "20"
            warning: "40"
            critical: "70"
```

### Disable this alert entirely

To completely disable DNSNxDomain alerts:

```bash
oc edit flowcollector cluster
```

Add DNSNxDomain to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - DNSNxDomain
```

For more information on configuring Network Observability alerts, see the
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

## Impact

High rates of NX_DOMAIN errors can indicate:

- Misconfigured applications trying to resolve non-existent domain names
- Typos in service URLs or DNS names in application configuration
- Suboptimal DNS searches
- Deleted services that applications are still trying to reach
- DNS-based service discovery issues
- Potential security concerns

While NX_DOMAIN errors are less critical than other DNS errors (like
SERVFAIL), a high rate can still cause:

- Wasted network resources and DNS server load
- Application delays due to failed lookups
- Logging noise and monitoring clutter
- Potential security blind spots if malicious activity is masked

NX_DOMAIN errors can be especially frequent in Kubernetes because of [DNS searches
on service and pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/),
which does not necessarily mean a misconfiguration or broken URL, but can still
negatively impact performance.

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

When NX_DOMAIN errors are returned despite of an apparently valid Service or Pod host name,
such as `my-svc.my-namespace.svc`, this is likely because the resolver is configured
to query DNS for different suffixes. It can be optimized by **adding a trailing dot on
fully-qualified domain names** which tells the resolver that the name is unambiguous.

For instance, instead of `https://my-svc.my-namespace.svc`, use:
`https://my-svc.my-namespace.svc.cluster.local.`, with trailing dot included.

For other mitigation strategies and solutions, refer to:

- [Troubleshooting DNS in OpenShift](https://access.redhat.com/solutions/3804501)
- [Kubernetes concepts: DNS for Services and Pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)
- [OpenShift DNS Operator](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/networking_operators/dns-operator)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
