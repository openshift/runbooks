# Ingress5xxErrors

## Meaning

The `Ingress5xxErrors` alert template is triggered when Network Observability
detects a high percentage of 5xx HTTP response codes from HAProxy ingress
traffic. This template can generate multiple alert instances depending on how
it's configured in the FlowCollector custom resource.

**Possible alert variants:**

- `Ingress5xxErrors_Critical` - Global cluster-wide ingress 5xx error rate
  exceeds critical threshold
- `Ingress5xxErrors_Warning` - Global cluster-wide ingress 5xx error rate
  exceeds warning threshold
- `Ingress5xxErrors_Info` - Global cluster-wide ingress 5xx error rate exceeds
  info threshold
- `Ingress5xxErrors_PerDstNamespace{Critical,Warning,Info}` - Ingress 5xx
  error rate for traffic destined to a specific namespace exceeds threshold

The alert fires when the percentage of 5xx errors in HAProxy ingress traffic
exceeds the configured threshold. By default, thresholds are set at 5% for
info level and 10% for warning level.

**Note:** This alert monitors HAProxy metrics, so it only tracks traffic
coming through OpenShift Routes and Ingress controllers. It does not track
internal cluster traffic or services exposed through LoadBalancer/NodePort.

### Switch to recording mode (alternative to alerts)

If you want to monitor ingress 5xx errors in the Network Health dashboard
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
      - template: Ingress5xxErrors
        mode: Recording
        variants:
        - groupBy: Namespace
          thresholds:
            info: "5"
            warning: "10"
            critical: "20"
```

In recording mode:

- Ingress 5xx error violations remain visible in the **Network Health**
  dashboard
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
      - template: Ingress5xxErrors
        mode: Alert
        variants:
        - groupBy: Namespace
          thresholds:
            info: "10"      
            warning: "20"   
            critical: "30"
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
