# KubePersistentVolumeInodesFillingUp

## Meaning

The `KubePersistentVolumeInodesFillingUp` alert triggers when a persistent
volume in one of the system namespaces has less than 3% of its allocated inodes
left. System namespaces include `default` and those that have names beginning
with `openshift-` or `kube-`.

## Impact

Significant inode usage by a system component is likely to prevent the
component from functioning normally. Signficant inode usage can also lead to a
partial or full cluster outage.

## Diagnosis

The alert labels include the name of the persistent volume claim (PVC)
associated with the volume running low on storage. The labels also include the
namespace in which the PVC is located. Use this information to graph
available storage in the OpenShift web console under Observe -> Metrics.  

The following is an example query for a PVC associated with a Prometheus
instance in the `openshift-monitoring` namespace:

```text
kubelet_volume_stats_inodes_used{
  namespace="openshift-monitoring",
  persistentvolumeclaim="prometheus-k8s-db-prometheus-k8s-0"
}
```

You can inspect the status of the volume manually to determine which directory
is consuming a large number of inodes:

```console
$ PVC_NAME='<persistentvolumeclaim label from alert>'
$ NAMESPACE='<namespace label from alert>'

$ oc -n $NAMESPACE describe pvc $PVC_NAME
$ POD_NAME='<"Used By:" field from the above output>'

$ oc -n $NAMESPACE rsh $POD_NAME
$ cd /path/to/pvc-mount
$ ls -li .
$ stat
```

## Mitigation

Mitigating this issue depends on the total count of files, directories, and
symbolic links. You cannot expand the number of inodes on a file system after
the file system has been created. However, you can adjust the configuration for
the component using the volume so that it creates fewer files, directories, and
symbolic links.