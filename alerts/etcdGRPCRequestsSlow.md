# etcdGRPCRequestsSlow

## Meaning

This alert fires when the 99th percentile of etcd gRPC requests are too slow.

## Impact

When requests are too slow, they can lead to various scenarios like leader
election failure, slow reads and writes.

## Diagnosis

This could be result of slow disk (due to fragmented state) or CPU contention.

### Slow disk

One of the most common reasons for slow gRPC requests is disk. Checking disk
related metrics and dashboards should provide a more clear picture.

#### PromQL queries used to troubleshoot

Verify the value of how slow the etcd gRPC requests are by using the following
query in the metrics console:

```console
histogram_quantile(0.99, sum(rate(grpc_server_handling_seconds_bucket{job=~".*etcd.*", grpc_type="unary"}[5m])) without(grpc_type))
```
That result should give a rough timeline of when the issue started.

`etcd_disk_wal_fsync_duration_seconds_bucket` reports the etcd disk fsync
duration, `etcd_server_leader_changes_seen_total` reports the leader changes. To
rule out a slow disk and confirm that the disk is reasonably fast, 99th
percentile of the etcd_disk_wal_fsync_duration_seconds_bucket should be less
than 10ms. Query in metrics UI:

```console
histogram_quantile(0.99, sum by (instance, le) (irate(etcd_disk_wal_fsync_duration_seconds_bucket{job="etcd"}[5m])))
```
#### Console dashboards

In the OpenShift dashboard console under Observe section, select the etcd
dashboard. There are both RPC rate as well as Disk Sync Duration dashboards
which will assit with further issues.


## Mitigation

### Fragmented state

In the case of slow fisk or when the etcd DB size increases, we can defragment
existing etcd DB to optimize DB consumption as described in
[here][etcdDefragmentation]. Run the following command in all etcd pods.

```console
$ etcdctl defrag
```

As validation, check the endpoint status of etcd members to know the reduced
size of etcd DB. Use for this purpose the same diagnostic approaches as listed
above. More space should be available now.

Further info on etcd best practices can be found in the [OpenShift docs
here][etcdPractices].

[etcdDefragmentation]: https://etcd.io/docs/v3.4.0/op-guide/maintenance/
[etcdPractices]: https://docs.openshift.com/container-platform/4.7/scalability_and_performance/recommended-host-practices.html#recommended-etcd-practices_
