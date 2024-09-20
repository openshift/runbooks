# HCOMisconfiguredDescheduler

## Meaning

A Descheduler is a OpenShift Container Platform application that causes the
control plane to
re-arrange the workloads in a better way.

The descheduler uses the OpenShift Container Platform eviction API to evict
pods, and receives
feedback from `kube-api` whether the eviction request was granted or not.
On the other side, in order to keep VM live and trigger live-migration,
OpenShift Virtualization handles eviction requests in a custom way and
unfortunately
a live migration takes time.
So from the descheduler's point of view, `virt-launcher` pods fail to be
evicted, but they actually migrating to another node in background.
So the way OpenShift Virtualization handles eviction requests causes the
descheduler to
make wrong decisions and take wrong actions that could destabilize the cluster.

To correctly handle the special case of `VM` pod evicted triggering a live
migration to another node, the `Kube Descheduler Operator` introduced
a `profileCustomizations` named `devEnableEvictionsInBackground`
which is currently considered an `alpha` feature
on `Kube Descheduler Operator` side.

## Impact

Using the descheduler operator for OpenShift Virtualization VMs without the
`devEnableEvictionsInBackground` profile customization can lead to
unstable or oscillatory behavior, undermining cluster stability.

## Diagnosis

1. Check the CR for `Kube Descheduler Operator`:

   ```bash
   $ oc get -n openshift-kube-descheduler-operator KubeDescheduler cluster -o yaml
   ```

looking for:

   ```yaml
   spec:
      profileCustomizations:
         devEnableEvictionsInBackground: true
   ```

If not there, the `Kube Descheduler Operator` is not correctly configured
to work alongside OpenShift Virtualization.

## Mitigation

Set:
   ```yaml
   spec:
      profileCustomizations:
         devEnableEvictionsInBackground: true
   ```
on the CR for the `Kube Descheduler Operator` or
remove the `Kube Descheduler Operator`.

Please notice that `EvictionsInBackground` is an alpha feature,
and it's subject to change and currently provided as an
experimental feature.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.