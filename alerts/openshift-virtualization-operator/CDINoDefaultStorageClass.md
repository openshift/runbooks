# CDINoDefaultStorageClass

## Meaning

This alert fires when there is no default (OpenShift Container Platform or
OpenShift Virtualization) storage
class, and a data volume is pending for one.

A default OpenShift Virtualization storage class has precedence over a default
OpenShift Container Platform
storage class for creating a VirtualMachine disk image.

## Impact

If there is no default OpenShift Container Platform or OpenShift Virtualization
storage class, a data volume that
does not have a specified storage class remains in a "pending" state.

## Diagnosis

1. Check for a default OpenShift Container Platform storage class by running
the following
command:

  ```bash
  $ oc get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubernetes.io/is-default-class"=="true")|.name'
  ```

2. Check for a default OpenShift Virtualization storage class by running the
following command:

  ```bash
  $ oc get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubevirt.io/is-default-virt-class"=="true")|.name'
  ```

## Mitigation

Create a default storage class for either OpenShift Container Platform or
OpenShift Virtualization or for both.

A default OpenShift Virtualization storage class has precedence over a default
OpenShift Container Platform
storage class for creating a virtual machine disk image.

* Create a default OpenShift Container Platform storage class by running the
following command:

  ```bash
  $ oc patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```

* Create a default OpenShift Virtualization storage class by running the
following command:

  ```bash
  $ oc patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubevirt.io/is-default-virt-class":"true"}}}'
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.