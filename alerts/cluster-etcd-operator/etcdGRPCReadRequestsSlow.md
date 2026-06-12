# etcdGRPCReadRequestsSlow

## Alert Description

| Impacted User   | Impact |
|:----------------|:-------|
| Customer Impact | High   |
| SREP Impact     | High   |

**Severity:** Critical

**Alert Expression:**
```promql
histogram_quantile(0.99, sum(rate(grpc_server_handling_seconds_bucket{job="etcd", grpc_method="Range", grpc_type="unary"}[10m])) without(grpc_type)) > 3
```

**Trigger Condition:**
The 99th percentile of etcd gRPC Range requests exceeds 3 seconds for 10 minutes.

**What it means:**
etcd is taking too long to serve critical read requests.
Range requests are the primary read operation in etcd's gRPC API, used by:
- `kubectl get` commands (listing pods, nodes, deployments, etc.)
- Kubernetes API server listing resources
- Watch operation initialization
- Any query for key-value data from etcd

**Common Root Causes:**
1. **Disk performance degradation** (~70% of cases) - Slow fsync operations
2. **CPU/Memory exhaustion** (~15% of cases) - Resource contention on control plane
3. **Network latency** (~10% of cases) - High peer round-trip time
4. **Database size/quota issues** (~5% of cases) - Large database or quota limit

## Troubleshooting

### Check for Correlated Critical Alerts

This alert's signal closely matches with some existing alerts.
Look for these commonly co-firing alerts for mitigation steps:

1. **etcdHighFsyncDurations (CRITICAL)** - Most likely cause - disk performance.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdHighFsyncDurations.md).

2. **etcdNoLeader (CRITICAL)** - May have fired 9+ minutes earlier.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdNoLeader.md).

3. **etcdInsufficientMembers (CRITICAL)** - Quorum loss scenario.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdInsufficientMembers.md).

4. **etcdHighNumberOfFailedGRPCRequests (CRITICAL)** - Requests timing out.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdHighNumberOfFailedGRPCRequests.md).

### Refer to the `etcd` team

If you are unsure about further steps to take, you can refer to common [etcd checks](https://github.com/openshift/ops-sop/blob/master/v4/troubleshoot/etcd.md#etcd-checks).
Collect information on current situation of the `etcd`.
You can also refer it to the `etcd` team in #forum-ocp-etcd.
