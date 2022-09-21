# etcdGRPCRequestsSlow

## Meaning

This alert fires when the 99th percentile of etcd gRPC requests are too slow.

## Impact

When requests are too slow, they can lead to various scenarios like leader
election failure, slow reads and writes and general cluster instability.

## Diagnosis

This could be result of slow disk, network or CPU starvation/contention.

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

When txn calls are slow, another culprit can be the network roundtrip between the nodes. You can observe this with:

```console
histogram_quantile(0.99, sum by (instance, le) (irate(etcd_network_peer_round_trip_time_seconds_bucket{job="etcd"}[5m])))
```


You can find more performance troubleshooting tips in [OpenShift etcd Performance Metrics](https://github.com/openshift/cluster-etcd-operator/blob/master/docs/performance-metrics.md).

#### Console dashboards

In the OpenShift dashboard console under Observe section, select the etcd
dashboard. There are both RPC rate as well as Disk Sync Duration dashboards
which will assist with further issues.

### Resource exhaustion

It can happen that etcd responds slower due to CPU resource exhaustion.
This was seen in some cases when one application was requesting too much CPU
which led to this alert firing for multiple methods.

Often if this is the case, we also see
`etcd_disk_wal_fsync_duration_seconds_bucket` slower as well.

To confirm this is the cause of the slow requests either:

1. In OpenShift console on primary page under "Cluster utilization" view the
   requested CPU vs available.

2. PromQL query is the following to see top consumers of CPU:

```console
      topk(25, sort_desc(
        sum by (namespace) (
          (
            sum(avg_over_time(pod:container_cpu_usage:sum{container="",pod!=""}[5m])) BY (namespace, pod)
            *
            on(pod,namespace) group_left(node) (node_namespace_pod:kube_pod_info:)
          )
          *
          on(node) group_left(role) (max by (node) (kube_node_role{role=~".+"}))
        )
      ))
```

### Rogue Workloads

In some cases, we've seen non-OpenShift workload put a lot of stress on the API server that eventually cascades into etcd. 
One specific instance was listing all pods across namespaces exhausting CPU and memory on API server and subsequently on etcd.

Please consult the audit log and see whether some service accounts make suspicious calls, both in terms of generality (listings, many/all namespaces) and frequency (eg listing all pods every 10s).


## Mitigation

Depending on what resource was determined to be exhausted, you can try the following:

### CPU

Find the offending process that uses too much CPU, try to limit or shutdown the process.
If feasible on clouds, adding more or faster CPUs may help to reduce the latency.


### Disk

Find the offending process that causes the disk performance to degrade, this can also be a noisy neighbour process on the control plane node (eg fluentd with logging, OVN) or etcd itself. If the culprit is determined to etcd, try to reduce the load coming from the apiserver. Most commonly this also happens when a cluster is scaled up with many more nodes, so reducing the cluster scale again can help.

If feasible on clouds, upgrading your storage or instance type can significantly increase your sequential IOPS and available bandwidth.


### Network

Ensure nothing is exhausting the network bandwidth as this causes package loss and increases the latency.
As with the previous two sections, try to isolate the offending process and mitigate from there. 
If etcd is the offender, try to reduce load or increase the available resources.

