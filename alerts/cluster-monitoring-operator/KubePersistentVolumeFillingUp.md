# KubePersistentVolumeFillingUp

## Meaning

This alert fires when a persistent volume in one of the system namespaces has
less than 3% of its total space remaining. System namespaces include
`default` and those that have names beginning with `openshift-` or `kube-`.
## Impact

If a persistent volume used by a system component fills up, the component
is unlikely to function normally. A full persistent volume can also lead to a
partial or full cluster outage.

## Diagnosis

The alert labels include the name of the persistent volume claim (PVC)
associated with the volume running low on storage. The labels also include the
namespace in which the PVC is located. Use this information to graph
available storage in the OpenShift web console under Observe -> Metrics.  

The following is an example query for a PVC associated with a Prometheus
instance in the `openshift-monitoring` namespace:

```text
kubelet_volume_stats_available_bytes{
  namespace="openshift-monitoring",
  persistentvolumeclaim="prometheus-k8s-db-prometheus-k8s-0"
}
```

You can also inspect the contents of the volume manually to determine what is
using the storage:

```console
$ PVC_NAME='<persistentvolumeclaim label from alert>'
$ NAMESPACE='<namespace label from alert>'

$ oc -n $NAMESPACE describe pvc $PVC_NAME
$ POD_NAME='<"Used By:" field from the above output>'

$ oc -n $NAMESPACE rsh $POD_NAME
$ df -h
```

## Mitigation

Mitigation for this issue depends on what is filling up the storage.  

You can try allocating more storage space to the affected volume to solve the
issue.

You can also try adjusting the configuration for the component that is the
volume so that the component requires less storage space. For example, if logs
for a component are filling up the persistent volume, you can change the log
level so that less information is logged and therefore less space is required
for logs.
