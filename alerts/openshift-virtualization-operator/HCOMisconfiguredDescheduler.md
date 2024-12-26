# HCOMisconfiguredDescheduler

## Meaning

A descheduler is a OpenShift Container Platform application that causes the
control plane to
re-arrange the workloads in a better way.

The descheduler uses the OpenShift Container Platform eviction API to evict
pods, and receives
feedback from `kube-api` on whether the eviction request was granted.
In contrast, to keep a VM live and trigger a live migration,
OpenShift Virtualization handles eviction requests in a custom manner,
and a live migration takes time to perform.

Therefore, when a `virt-launcher` pod is migrating to another node in the
background,
the descheduler detects this as a pod that failed to be evicted. As a
consequence,
the manner in which OpenShift Virtualization handles eviction requests causes
the descheduler
to make incorrect decisions and take incorrect actions that might
destabilize the cluster.

To correctly handle the special case of an evicted VM pod triggering a live
migration to another node, the `Kube Descheduler Operator` introduced
a `profileCustomizations` named `devEnableEvictionsInBackground`.
This is currently considered an `alpha` feature for
on `Kube Descheduler Operator`.

## Impact

Using the descheduler operator for KubeVirt VMs without the
`devEnableEvictionsInBackground` profile customization might lead
to unstable or unpredictable behavior, which negatively impacts cluster
stability.

## Diagnosis

1. Check the CR for `Kube Descheduler Operator`:

   ```bash
   $ oc get -n openshift-kube-descheduler-operator KubeDescheduler cluster -o yaml
   ```

2. Search for the following lines in the CR:

   ```yaml
   spec:
      profileCustomizations:
         devEnableEvictionsInBackground: true
   ```

If these lines are not present, the `Kube Descheduler Operator` is not correctly
configured to work alongside OpenShift Virtualization.

## Mitigation

Do one of the following:

* Add the following lines to the CR for `Kube Descheduler Operator`:
   ```yaml
   spec:
      profileCustomizations:
         devEnableEvictionsInBackground: true
   ```

* Remove the `Kube Descheduler Operator`.

Note that `EvictionsInBackground` is an alpha feature,
and as such, it is provided as an experimental feature
and is subject to change.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.