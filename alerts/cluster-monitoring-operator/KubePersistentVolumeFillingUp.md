# KubePersistentVolumeFillingUp

## Meaning

This alert fires when a persistent volume in one of the system namespaces,
i.e. a namespace beginning with `openshift-`, `kube-`, or the `default`
namespace, has less than 3% of its total space left.

## Impact

A full persistent volume used by a system component is likely to prevent the
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
kubelet_volume_stats_available_bytes{
  namespace="openshift-monitoring",
  persistentvolumeclaim="prometheus-k8s-db-prometheus-k8s-0"
}
```

You can inspect the contents of the volume manually to determine what is using
the storage:

```console
$ PVC_NAME='<persistentvolumeclaim label from alert>'
$ NAMESPACE='<namespace label from alert>'

$ oc -n $NAMESPACE describe pvc $PVC_NAME
$ POD_NAME='<"Used By:" field from the above output>'

$ oc -n $NAMESPACE rsh $POD_NAME
$ df -h
```

## Mitigation

The mitigation largely depends on what is using the storage.  You may be able to
simply allocate more storage to the affected volume, or adjust the configuration
for the component using the volume to use less space -- for instance by lowering
a logging level or similar.
