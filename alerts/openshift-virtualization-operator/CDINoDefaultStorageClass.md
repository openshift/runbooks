# CDINoDefaultStorageClass

## Meaning

This alert fires when no default OpenShift Container Platform or
OpenShift Virtualization storage class is defined.

## Impact

If no default OpenShift Container Platform or OpenShift Virtualization storage
class is defined, a data volume requesting a default storage class (the storage
class is not specified), remains in a "pending" state.

## Diagnosis

1. Check for a default OpenShift Container Platform storage class by running
the following command:

   ```bash
   $ oc get sc -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubevirt\.io/is-default-class=="true")].metadata.name}'
   ```

2. Check for a default OpenShift Virtualization storage class by running
the following command:

   ```bash
   $ oc get sc -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubevirt\.io/is-default-virt-class=="true")].metadata.name}'
   ```

## Mitigation

Create a default storage class for either OpenShift Container Platform or
OpenShift Virtualization or for both.

A default OpenShift Virtualization storage class has precedence over a default
OpenShift Container Platform storage class for creating a virtual machine disk image.

* Create a default OpenShift Container Platform storage class by running
  the following command:

  ```bash
  $ oc patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```

* Create a default OpenShift Virtualization storage class by running
  the following command:

  ```bash
  $ oc patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubevirt.io/is-default-virt-class":"true"}}}'
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
