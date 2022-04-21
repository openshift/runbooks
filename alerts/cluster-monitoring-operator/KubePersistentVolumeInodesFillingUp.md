# KubePersistentVolumeInodesFillingUp

## Meaning

This alert fires when a persistent volume in one of the system namespaces,
i.e. a namespace beginning with `openshift-`, `kube-`, or the `default`
namespace, has less than 3% of its allocated inodes left.

## Impact

A significant inode usage by a system component is likely to prevent the
component from functioning normally, and may lead to a partial or full cluster
outage.

## Diagnosis

The alert labels should include the name of the PersistentVolumeClaim associated
with the volume that is low on storage, as well as the namespace that claim is
in.  You can use these to graph the available storage in the OpenShift web
console under Observer -> Metrics.  The following is an example query for a
volume claim associated with a Prometheus instance in the `openshift-monitoring`
namespace:

```text
kubelet_volume_stats_inodes_used{
  namespace="openshift-monitoring",
  persistentvolumeclaim="prometheus-k8s-db-prometheus-k8s-0"
}
```

You can inspect the status of the volume manually to determine which directory
consumes large number of the inode:

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

The mitigation largely depends on the count of files/directories/soft links.
It is not possible to expand the number of inodes on a filesystem after
it is created, better adjust the configuration
for the component using the volume to create less number of
files/directories/soft links.
