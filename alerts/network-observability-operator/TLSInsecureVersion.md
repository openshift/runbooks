# TLSInsecureVersion

## Meaning

The `TLSInsecureVersion` health rule template is triggered when Network Observability
detects a high percentage of TLS traffic using insecure protocol versions (TLS 1.0,
TLS 1.1, SSL 2.0, or SSL 3.0). These legacy protocols have known security
vulnerabilities and should not be used in production environments.

The rule can generate multiple alert or recording instances depending on how it's
configured in the `FlowCollector` custom resource. Both the Alert and the Recording
modes are displayed in the Network Health view, but only the Alert mode can
generate Prometheus alerts.

**Note:** This rule requires the `TLSTracking` agent feature to be enabled
in the `FlowCollector` configuration.

**Possible alert variants:**

- `TLSInsecureVersion_Critical` - Global cluster-wide insecure TLS usage exceeds
  critical threshold
- `TLSInsecureVersion_Warning` - Global cluster-wide insecure TLS usage exceeds
  warning threshold
- `TLSInsecureVersion_Info` - Global cluster-wide insecure TLS usage exceeds
  info threshold
- `TLSInsecureVersion_PerDstNamespace{Critical,Warning,Info}` - Insecure TLS usage
  for traffic destined to a specific namespace exceeds threshold
- `TLSInsecureVersion_PerSrcNamespace{Critical,Warning,Info}` - Insecure TLS usage
  for traffic originating from a specific namespace exceeds threshold
- `TLSInsecureVersion_PerDstNode{Critical,Warning,Info}` - Insecure TLS usage for
  traffic destined to a specific node exceeds threshold
- `TLSInsecureVersion_PerSrcNode{Critical,Warning,Info}` - Insecure TLS usage for
  traffic originating from a specific node exceeds threshold
- `TLSInsecureVersion_PerDstWorkload{Critical,Warning,Info}` - Insecure TLS usage
  for traffic destined to a specific workload exceeds threshold
- `TLSInsecureVersion_PerSrcWorkload{Critical,Warning,Info}` - Insecure TLS usage
  for traffic originating from a specific workload exceeds threshold

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
      - template: TLSInsecureVersion
        mode: Alert
        variants:
        - thresholds:
            warning: "5"
          groupBy: Namespace
```

If you prefer to switch to the recording mode with `mode: Recording`:

- TLS insecure version violations remain visible in the **Network Health** dashboard
- No Prometheus alerts are generated
- Metrics are still calculated and stored as recording rules
- Useful for teams that prefer passive monitoring without alert fatigue

### Disable this alert entirely

To completely disable TLSInsecureVersion alerts:

```bash
oc edit flowcollector cluster
```

Add TLSInsecureVersion to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - TLSInsecureVersion
```

For more information on configuring Network Observability alerts, see the
[Network Observability documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/network-observability-alerts_nw-observe-network-traffic).

## Impact

Using insecure TLS versions exposes communications to serious security risks:

- Man-in-the-middle attacks that can intercept and modify traffic
- Protocol downgrade attacks that force use of weaker encryption
- Cipher suite vulnerabilities such as BEAST, POODLE, and others
- Compliance violations for standards like PCI DSS, HIPAA, and SOC 2
- Data breaches and unauthorized access to sensitive information

Modern security standards mandate TLS 1.2 or higher. Even a small percentage
of traffic using insecure TLS versions can represent a significant security
risk, as attackers can exploit these connections to gain unauthorized access.

Industry requirements include:
- PCI DSS 3.2.1 requires TLS 1.2+ (deprecated TLS 1.0/1.1 since June 2018)
- NIST SP 800-52 Rev. 2 prohibits TLS 1.0 and 1.1
- Most modern browsers have disabled TLS 1.0/1.1 by default

## Diagnosis

When this alert fires, you can investigate further by using the Network
Observability interface:

1. **Navigate to alert details**: Click on the alert in the Network Health
   dashboard to view specific details of the alert, including which namespaces,
   nodes, or workloads are using insecure TLS versions.

2. **Navigate to network traffic**: From the alert details, you can navigate
   to the Network Traffic view to examine the specific flows that are related
   to this alert. This allows you to see:
   - Source and destination of the TLS traffic
   - Detailed flow information including TLS version used
   - Which applications are acting as clients vs servers

For additional troubleshooting resources, refer to the documentation links in
the Mitigation section below.

## Mitigation

For mitigation strategies and solutions, refer to:

- [OpenShift TLS Security Profiles](https://docs.openshift.com/container-platform/latest/security/tls-security-profiles.html)
- [Configuring ingress cluster traffic](https://docs.openshift.com/container-platform/latest/networking/configuring_ingress_cluster_traffic/overview-traffic.html)
- [NIST SP 800-52 Rev. 2 - Guidelines for TLS](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-52r2.pdf)
- [OpenShift Networking](https://docs.redhat.com/en/documentation/openshift_container_platform/latest#Networking)
