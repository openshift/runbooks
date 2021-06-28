# HighlyAvailableWorkloadIncorrectlySpread

## Meaning

Highly available workloads with persistent storage are incorrectly spread across
multiple nodes. In high availability topology, all workloads with multiple
replicas should be spread across multiple nodes to avoid having a single point
of failure. However, workloads with persistent storage are more complex to
schedule since persistent volumes may be bound to an availability zone, so they
need to be spread manually before any scheduling constraints are put in place.

## Impact

The cluster isn't fully highly available since workloads with persistent
storage might still have a single point of failure.

## Diagnosis

The alert indicates which workloads aren't correctly spread across multiple
nodes. You can find information to identify them under the `workload` and
`namespace` labels.

```console
 - alertname = HighlyAvailableWorkloadIncorrectlySpread
...
 - workload = prometheus-k8s
...
 - namespace = openshift-monitoring
```

Then you can verify that multiple instances of the same application are indeed
scheduled on the same node by running:

```console
$ oc -n "$NS" get -o wide pods | grep "$WORKLOAD"
```

## Mitigation

To mitigate this issue, you can follow the steps below, but first you should
get the following information from the alert labels:
* `workload`: workload name
* `namespace`: namespace of the workload
* `node`: node hosting multiple instances of the workload

You should start by cordoning the node to prevent pods from being rescheduled
on it:

```console
$ oc adm cordon "$NODE"
```

Then you will need to get the list of pods that should be rescheduled for the
given workload and node:

```console
$ oc -n "$NS" get -o wide pods | grep "$WORKLOAD.*$NODE" | cut -f1 -d ' ' | head -n-1)
```

If your storage system is bound to availability zones, you will also want to
get the number of availability zones present on your cluster to know if you
need to delete persistent volume claims:

```console
$ oc get nodes -o yaml | grep -E "^\s+failure-domain.beta.kubernetes.io/zone" | uniq | wc -l
```

Then for each pod, you will want to delete the pod and the PVC attached to it,
if required by the previous step, in order to reschedule it on another node:

> Note that deleting the PVC will result in data loss and the replication level
will decrease, but this is a required step if you want your cluster to be
highly available.

```console
$ PVC=$(oc get -n "$NS" pod "$POD" -ojson | jq -r '.spec.volumes[] | select(.persistentVolumeClaim!=null) | .persistentVolumeClaim.claimName')
$ oc delete -n "$NS" pvc "$PVC"
$ oc delete -n "$NS" pod "$POD"
```

Once all the pods are rescheduled, you can uncordon the node:

```console
$ oc adm uncordon "$NODE"
```
