# CDINoDefaultStorageClass

## Meaning

This alert fires when a data volume is `Pending` because there is no default
storage class.

A default virtualization storage class has precedence over a default OpenShift
Container Platform
storage class for creating a VirtualMachine disk image.

## Impact

If there is no default OpenShift Container Platform storage class and no
default virtualization
storage class, a data volume that does not have a specified storage class
remains in a `Pending` phase.

## Diagnosis

1. Check for a default OpenShift Container Platform storage class by running
the following
command:

  ```bash
  $ oc get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubernetes\.io\/is-default-class=="true")].metadata.name}'
  ```

2. Check for a default virtualization storage class by running the following
command:

  ```bash
  $ oc get sc -o jsonpath='{.items[?(.metadata.annotations.storageclass\.kubevirt\.io\/is-default-virt-class=="true")].metadata.name}'
  ```

## Mitigation

Create a default storage class for OpenShift Container Platform,
virtualization, or both.

A default virtualization storage class has precedence over a default OpenShift
Container Platform
storage class for creating a virtual machine disk image.

* Create a default OpenShift Container Platform storage class by running the
following command:

  ```bash
  $ oc patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```

* Create a default virtualization storage class by running the following
command:

  ```bash
  $ oc patch storageclass <storage-class-name> -p '{"metadata": {"annotations":{"storageclass.kubevirt.io/is-default-virt-class":"true"}}}'
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.