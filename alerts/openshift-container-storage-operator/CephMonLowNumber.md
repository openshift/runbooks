# CephMonLowNumber

## Meaning

The number of ceph monitors in the cluster can be adjusted to improve cluster
resiliency.
Typically the number of failure zones in the cluster is 5 or more, and there
are only 3 monitors.

## Impact

This a "info" level alert, and therefore just a suggestion.
The alert is just suggesting to increase the number of ceph monitors, to be
more resistent to failures.
It can be silenced without any impact in the cluster functionality or
performance.
If the number of monitors is increased to 5, the cluster will be more robust.

## Diagnosis

Check the number of Ceph Monitors:

```bash
    oc get pods -l app=rook-ceph-mon --no-headers=true -n openshift-storage | wc -l
```

Check the number of failure zones available:

```bash
    oc get nodes -o jsonpath='{.items[*].metadata.labels.topology\.kubernetes\.io/zone}' | tr ' ' '\n' | sort -u | wc -l
```

It the number of available failure zones is greater or equal to 5, and there
are only 3 monitors, the alert will be raised.

## Mitigation

If increasing the number of monitors to 5 is not a right option (for any cause),
the alert can be silenced.

If to increase the number of monitors is an acceptable proposal, then execute
the following command to do that:

```bash
    oc patch storageclusters.ocs.openshift.io ocs-storagecluster -n openshift-storage --type merge --patch '{"spec": {"managedResources": {"cephCluster": {"monCount" : 5}}}}'
```

After scaling the deployment, monitor the creation and readiness of new monitor
 pods using:

```bash
    oc get pods -n openshift-storage -l app=rook-ceph-mon
```
