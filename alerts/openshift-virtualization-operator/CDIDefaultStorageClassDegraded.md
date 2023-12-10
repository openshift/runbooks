# CDIDefaultStorageClassDegraded

## Meaning

This alert fires when there is no default storage class that supports smart cloning
(CSI or snapshot-based) or the ReadWriteMany access mode.

## Impact

If the default storage class does not support smart cloning, the default cloning
method is host-assisted cloning, which is much less efficient.

If the default storage class does not support ReadWriteMany, virtual machines (VMs)
cannot be live migrated.

Note: A default OpenShift Virtualization storage class has precedence over a
default OpenShift Container Platform storage class when creating a
VM disk.

## Diagnosis

1. Get the default OpenShift Virtualization storage class by running the following
command:

   ```bash
   $ oc get sc -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubevirt\.io/is-default-virt-class=="true")].metadata.name}'
   ```

2. If a default OpenShift Virtualization storage class exists, check that it
supports ReadWriteMany by running the following command:

   ```bash
   $ oc get storageprofile <storage_class> -o json | jq '.status.claimPropertySets'| grep ReadWriteMany
   ```

3. If there is no default OpenShift Virtualization storage class, get the
default OpenShift Container Platform storage class by running the following
command:

   ```bash
   $ oc get sc -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubevirt\.io/is-default-class=="true")].metadata.name}'
   ```

4. If a default OpenShift Container Platform storage class exists, check that it
supports ReadWriteMany by running the following command:

   ```bash
   $ oc get storageprofile <storage_class> -o json | jq '.status.claimPropertySets'| grep ReadWriteMany
   ```

## Mitigation

Ensure that you have a default storage class, either OpenShift Container Platform
or OpenShift Virtualization, and that the default storage class supports
smart cloning and ReadWriteMany.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case, attaching
the artifacts gathered during the diagnosis procedure.