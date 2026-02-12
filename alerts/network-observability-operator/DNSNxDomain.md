# DNSNxDomain

## Meaning

The `DNSNxDomain` health rule template is triggered when Network Observability
detects a high percentage of DNS NX_DOMAIN errors. NX_DOMAIN indicates that
the queried domain name does not exist, potentially because of ambiguous DNS
searches.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule requires the `DNSTracking` agent feature to be enabled
in the `FlowCollector` configuration.

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
      - template: DNSNxDomain
        mode: Recording
        variants:
        - thresholds:
            info: "10"
            warning: "80"
          groupBy: Namespace
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

When NX_DOMAIN errors are returned despite of an apparently valid Service
or Pod host name, such as `my-svc.my-namespace.svc`, this is likely
because the resolver is configured to query DNS for different suffixes. It
can be optimized by **adding a trailing dot on fully-qualified domain
names** which tells the resolver that the name is unambiguous.

For instance, instead of `https://my-svc.my-namespace.svc`, use:
`https://my-svc.my-namespace.svc.cluster.local.`, with trailing dot included.

For other mitigation strategies and solutions, refer to:

- [Troubleshooting DNS in OpenShift](https://access.redhat.com/solutions/3804501)
- [Kubernetes concepts: DNS for Services and Pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)
- [OpenShift DNS Operator](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/networking_operators/dns-operator)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
