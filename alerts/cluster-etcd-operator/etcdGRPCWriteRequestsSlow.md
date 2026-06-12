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

**Trigger Condition:** The 99th percentile of etcd gRPC Txn (transaction/write) requests exceeds 5 seconds for 10 consecutive minutes.

**What it means:** etcd is taking too long to commit write requests, which are critical for maintaining cluster state. Transaction (Txn) requests are the primary write operation in etcd's gRPC API, used by:
- `kubectl apply/create/delete/patch` commands
- Kubernetes API server writing any resource changes
- Pod scheduling and lifecycle operations
- Resource updates by controllers and operators
- Any create/update/delete operation in the cluster
- Cluster state modifications requiring Raft consensus

## Troubleshooting

### Check for Correlated Critical Alerts

This is a new alert introduced in OCP 4.22. It's signal closely matches with some existing alerts. Look for these commonly co-firing alerts to mitigate the issue using known resolution steps:

1. **etcdHighFsyncDurations (CRITICAL)** - Most likely cause (95% correlation) - disk write performance. Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdHighFsyncDurations.md).
3. **etcdHighNumberOfFailedGRPCRequests (CRITICAL)** - Write requests timing out/failing (75% correlation). Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdHighNumberOfFailedGRPCRequests.md).
4. **etcdGRPCReadRequestsSlow (CRITICAL)** - Same root cause affects reads too (65% correlation). Refer to the SOP [here](https://github.com/openshift/ops-sop/blob/master/v4/alerts/etcdGRPCReadRequestsSlow.md).

### Refer to the `etcd` team

If you are unsure about further steps to take, you can refer to common [etcd checks](https://github.com/openshift/ops-sop/blob/master/v4/troubleshoot/etcd.md#etcd-checks) to collect information on current situation of the `etcd` and refer it to the `etcd` team in #forum-ocp-etcd.
