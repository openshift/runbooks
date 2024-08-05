# CDIDefaultStorageClassDegraded

## Meaning

This alert fires if the default storage class does not support smart cloning
(CSI or snapshot-based) or the ReadWriteMany access mode. The alert does not
fire if at least one default storage class supports these features.

A default virtualization storage class has precedence over a default OpenShift
Container Platform
storage class for creating a VirtualMachine disk image.

In case of single-node OpenShift, the alert is suppressed if there is a default
storage class that supports smart cloning, but not ReadWriteMany.

## Impact

If the default storage class does not support smart cloning, the default cloning
method is host-assisted cloning, which is much less efficient.

If the default storage class does not support ReadWriteMany, virtual machines
(VMs) cannot be live migrated.

## Diagnosis

1. Get the default virtualization storage class by running the following
command:

   ```bash
   $ export CDI_DEFAULT_VIRT_SC="$(oc get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubernetes\.io\/is-default-class=="true")].metadata.name}')"
   ```

2. If a default virtualization storage class exists, check that it supports
ReadWriteMany by running the following command:

   ```bash
   $ oc get storageprofile $CDI_DEFAULT_VIRT_SC -o jsonpath='{.status.claimPropertySets}' | grep ReadWriteMany
   ```

3. If there is no default virtualization storage class, get the default
OpenShift Container Platform storage class by running the following command:

   ```bash
   $ export CDI_DEFAULT_K8S_SC="$(oc get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubernetes\.io\/is-default-class=="true")].metadata.name}')"
   ```

4. If a default OpenShift Container Platform storage class exists, check that
it supports
ReadWriteMany by running the following command:

   ```bash
   $ oc get storageprofile $CDI_DEFAULT_VIRT_SC -o jsonpath='{.status.claimPropertySets}' | grep ReadWriteMany
   ```

## Mitigation

Ensure that you have a default (OpenShift Container Platform or virtualization)
storage class, and
that the default storage class supports smart cloning and ReadWriteMany.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case, attaching
the artifacts gathered during the diagnosis procedure.