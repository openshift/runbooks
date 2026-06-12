# etcdGRPCWriteRequestsSlow

## Alert Description

| Impacted User   | Impact |
|:----------------|:-------|
| Customer Impact | High   |
| SREP Impact     | High   |

**Severity:** Critical

**Alert Expression:**
```promql
histogram_quantile(0.99, sum(rate(grpc_server_handling_seconds_bucket{job="etcd", grpc_method="Txn", grpc_type="unary"}[10m])) without(grpc_type)) > 5
```

**Trigger Condition:**
The 99th percentile of etcd gRPC Txn (write) requests exceeds 5 seconds for 10 minutes.

**What it means:**
etcd is taking too long to commit critical write requests.
Transaction requests are the primary write operation in etcd's gRPC API, used by:

- `kubectl apply/create/delete/patch` commands
- Kubernetes API server writing any resource changes
- Pod scheduling and lifecycle operations
- Resource updates by controllers and operators
- Any create/update/delete operation in the cluster
- Cluster state modifications requiring Raft consensus

## Troubleshooting

### Check for Correlated Critical Alerts

This alert's signal closely matches with some existing alerts.
Look for these commonly co-firing alerts for mitigation steps:

1. **etcdHighFsyncDurations (CRITICAL)** - Most likely cause- disk write performance.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdHighFsyncDurations.md).

3. **etcdHighNumberOfFailedGRPCRequests (CRITICAL)** - Write requests timing out/failing.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdHighNumberOfFailedGRPCRequests.md).

4. **etcdGRPCReadRequestsSlow (CRITICAL)** - Same root cause affects reads too.
Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdGRPCReadRequestsSlow.md).

### Refer to the `etcd` team

Collect information on current situation of the `etcd`.
You can also refer it to the `etcd` team in #forum-ocp-etcd.
