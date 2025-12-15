# etcdDatabaseQuotaLowSpace

## Meaning

This alert fires when the total existing DB size exceeds 95% of the maximum
DB quota. The consumed space is in Prometheus represented by the metric
`etcd_mvcc_db_total_size_in_bytes`, and the DB quota size is defined by
`etcd_server_quota_backend_bytes`.

## Impact

In case the DB size exceeds the DB quota, no writes can be performed anymore on
the etcd cluster. This further prevents any updates in the cluster, such as the
creation of pods.

## Diagnosis

### Automated Analysis (Non-HCP Clusters)

Configuration Anomaly Detection (CAD) automatically analyzes this alert for
non-HCP clusters and creates a backplane cluster report that includes:

- Top Space Consumers by Namespace
- Largest ConfigMaps and Secrets
- Event storage by namespace

**Note**: Check the PagerDuty incident notes for the command to view the
cluster report. If the notes report a failure to clean up the snapshot, you
must access the node and manually clean up the etcd snapshot. The target node
and file path are provided in the PagerDuty incident notes.

#### Viewing Cluster Reports

List all cluster reports for the cluster:
```console
$ osdctl cluster reports list --cluster-id <$CLUSTER_ID>
```

View a specific report:
```console
$ osdctl cluster reports get --cluster-id <$CLUSTER_ID> --report-id <$REPORT_ID>
```

### CLI Checks

To run `etcdctl` commands, we need to `rsh` into the `etcdctl` container of any
etcd pod.

```console
$ oc rsh -c etcdctl -n openshift-etcd $(oc get pod -l app=etcd -oname -n openshift-etcd | awk -F"/" 'NR==1{ print $2 }')
```

Validate that the `etcdctl` command is available:

```console
$ etcdctl version
```

`etcdctl` can be used to fetch the DB size of the etcd endpoints.

```console
$ etcdctl endpoint status -w table
```

### PromQL queries

Check the percentage consumption of etcd DB with the following query in the
metrics console:

```console
(etcd_mvcc_db_total_size_in_bytes / etcd_server_quota_backend_bytes) * 100
```

Check the DB size in MB that can be reduced after defragmentation:

```console
(etcd_mvcc_db_total_size_in_bytes - etcd_mvcc_db_total_size_in_use_in_bytes)/1024/1024
```

## Mitigation

### Capacity planning

If the `etcd_mvcc_db_total_size_in_bytes` shows that you are growing close to
the `etcd_server_quota_backend_bytes`, etcd almost reached max capacity and it's
start planning for new cluster.

In the meantime before migration happens, you can use defrag to gain some time.

### Defrag

When the etcd DB size increases, we can defragment existing etcd DB to optimize
DB consumption as described in [here][etcdDefragmentation]. Run the following
command in all etcd pods.

```console
$ etcdctl defrag
```

As validation, check the endpoint status of etcd members to know the reduced
size of etcd DB. Use for this purpose the same diagnostic approaches as listed
above. More space should be available now.

[etcdDefragmentation]: https://etcd.io/docs/v3.4.0/op-guide/maintenance/
