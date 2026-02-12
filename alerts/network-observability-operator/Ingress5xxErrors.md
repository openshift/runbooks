# Ingress5xxErrors

## Meaning

The `Ingress5xxErrors` health rule template is triggered when Network Observability
detects a high percentage of 5xx HTTP response codes from HAProxy ingress
traffic.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generates Prometheus alerts.

**Note:** This rule monitors HAProxy metrics, so it only tracks traffic
coming through OpenShift Routes and Ingress controllers. It does not track
internal cluster traffic or services exposed through LoadBalancer/NodePort.

**Possible alert variants:**

- `Ingress5xxErrors_Critical` - Global cluster-wide ingress 5xx error rate
  exceeds critical threshold
- `Ingress5xxErrors_Warning` - Global cluster-wide ingress 5xx error rate
  exceeds warning threshold
- `Ingress5xxErrors_Info` - Global cluster-wide ingress 5xx error rate exceeds
  info threshold
- `Ingress5xxErrors_PerDstNamespace{Critical,Warning,Info}` - Ingress 5xx
  error rate for traffic destined to a specific namespace exceeds threshold

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
      - template: Ingress5xxErrors
        mode: Recording
        variants:
        - thresholds:
            info: "5"
            warning: "10"
          groupBy: Namespace
```

### Disable this alert entirely

To completely disable Ingress5xxErrors alerts:

```bash
oc edit flowcollector cluster
```

Add Ingress5xxErrors to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - Ingress5xxErrors
```

For more information on configuring Network Observability alerts, see the
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

## Impact

High rates of 5xx errors in ingress traffic indicate server-side problems that
can severely impact user experience and application availability:

- **Application failures**: Backend services are crashing, timing out, or
  returning errors
- **Database issues**: Connection pool exhaustion or database failures
- **Resource exhaustion**: Pods running out of memory, CPU, or reaching
  connection limits
- **Configuration errors**: Misconfigured applications or dependencies
- **Cascading failures**: One failing service causing failures in dependent
  services

5xx errors are particularly critical because they indicate problems within
your application or cluster infrastructure, not client-side issues (which
would generate 4xx errors).

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

Possible causes and mitigations include:

- Checking application logs and recent deployments for errors or changes
- Verifying backend pod health, readiness probes, and resource limits
- Scaling backend pods if they're overloaded
- Reviewing HAProxy router logs and metrics
- Checking for network policies blocking traffic

For mitigation strategies and solutions, refer to:

- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
