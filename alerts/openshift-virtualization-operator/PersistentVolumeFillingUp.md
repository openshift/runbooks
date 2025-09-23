<!--
Original source from KubePersistentVolumeFillingUp

https://github.com/prometheus-operator/runbooks/blob/main/content/runbooks/kubernetes/KubePersistentVolumeFillingUp.md
-->

# PersistentVolumeFillingUp

## Meaning

A persistent volume associated with your cluster has crossed the warning
threshold of 90% used capacity and is expected to fill up within 4 days.

## Impact

Service degradation, switching to read-only mode.

## Diagnosis

- Check persistent volume claim (PVC) usage from within a pod that mounts the
PVC. You can use the *exec* command to check the mounted filesystem in a pod:

```shell
$ oc exec -it <pod-name> -n <namespace> -- df -h /path/to/mount
```

- Check if services are enabled that might use large amounts of storage.
For example, these might include snapshots or automatic data retention.

## Mitigation

### Data retention

Deleting data that you no longer need is the fastest and the cheapest solution.

Request the service owner to delete specific sets of old data.

### Data export

If data is not needed in the service but needs to be processed later,
move it to a different storage resource, such as an S3 bucket.

### Data rebalance in the cluster

Some services automatically rebalance data on the cluster when a node fills up.
Some of these make it possible to rebalance data across existing nodes, others
may require adding new nodes. If data rebalancing is supported in your cluster,
increase the number of replicas and wait for data migration, or trigger the
migration manually.

Example services that support data rebalancing:

- cassandra
- ceph
- elasticsearch/opensearch
- gluster
- hadoop
- kafka
- minio

**Note**: Some services may require special scaling conditions for data
rebalancing, such as doubling the number of active nodes.

### Direct Volume resizing

If volume resizing is available, you can increase the capacity of the volume:

1. Check that volume expansion is available. To do so, use the following command
and replace `<my-namespace>` and `<my-pvc>`:

   ```shell
   $ oc get storageclass `oc -n <my-namespace> get pvc <my-pvc> -ojson | jq -r '.spec.storageClassName'`
   NAME                 PROVISIONER            RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
   standard (default)   kubernetes.io/gce-pd   Delete          Immediate           true                   28d
   ```

   This example displays  `ALLOWVOLUMEEXPANSION` as `true`, which means you can
   use volume resizing.

2. Resize the volume:
   ```shell
   $ oc -n <my-namespace> edit pvc <my-pvc>
   ```

3. Edit `.spec.resources.requests.storage` to use the required storage size.
If this succeeds, the PVC status will state "Waiting for user to (re-)start a
pod to finish file system resize of volume on node."

4. Verify that the setting has been changed accordingly:
    ```shell
    $ oc -n <my-namespace> get pvc <my-pvc>
    ```

5. When prompted by the PVC status, restart the respective pod. The following
command automatically finds the pod that mounts the PVC and deletes it. However,
if you know the pod name, you can also directly delete that pod:

   ```shell
   $ oc -n <my-namespace> delete pod `oc -n <my-namespace> get pod -ojson | jq -r '.items[] | select(.spec.volumes[] .persistentVolumeClaim.claimName=="<my-pvc>") | .metadata.name'`
   ```

### Migrate data to a new, larger volume

If resizing is not available and the data is not safe to delete, the best
solution is to create a larger volume and migrate the data.

### Purge the volume

If the data is ephemeral and volume expansion is not available, it may be best
to purge the volume.

**WARNING**: This will permanently delete the data on the volume.

### Migrate data to a new, larger instance pool in the same cluster

In very specific scenarios, it is better to schedule data migration in the same
cluster, but to a new instance. This may be difficult to accomplish due to how
certain resources are managed in OpenShift Container Platform.

The general procedure is as follows:
1. Add new nodes with greater capacity than the existing cluster.
2. Trigger data migration.
3. Scale the original instance pool to 0, and then delete it.

### Migrate data to a new, larger cluster

This is the most common scenario for resolving filling up persistent volumes.
However, it is significantly more expensive and may be time-consuming.
Also, migrating to a new cluster might cause split-brain issues.
As an example, the general procedure may have the following steps:

1. Create a data snapshot on the existing cluster.
2. Add a new cluster with greater capacity than the existing cluster.
3. Start a data restore operation on the new cluster based on the snapshot.
4. Switch the old cluster to read-only mode
5. Reconfigure networking to point to the new cluster.
6. Trigger data migration from the old cluster to the new cluster to sync the
differences between the snapshot and the latest writes.
7. Remove the old cluster.

### Support

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.